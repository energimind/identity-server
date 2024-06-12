package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/auth"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
	"github.com/energimind/identity-server/internal/core/infra/oauth/providers"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
)

const sessionTTL = 24 * time.Hour

// Service manages user sessions.
//
// It implements the auth.Service interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type Service struct {
	applicationFinder admin.ApplicationLookupService
	providerFinder    admin.ProviderLookupService
	apiKeyFinder      admin.APIKeyLookupService
	idGenerator       domain.IDGenerator
	sessionCache      domain.Cache
}

// NewService returns a new Service instance.
func NewService(
	applicationFinder admin.ApplicationLookupService,
	providerFinder admin.ProviderLookupService,
	apiKeyFinder admin.APIKeyLookupService,
	idgen domain.IDGenerator,
	cache domain.Cache,
) *Service {
	return &Service{
		applicationFinder: applicationFinder,
		providerFinder:    providerFinder,
		apiKeyFinder:      apiKeyFinder,
		idGenerator:       idgen,
		sessionCache:      cache,
	}
}

// Ensure service implements the auth.Service interface.
var _ auth.Service = (*Service)(nil)

// ProviderLink implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) ProviderLink(ctx context.Context, applicationCode, providerCode, action string) (string, error) {
	const defaultAction = "login"

	app, err := s.applicationFinder.LookupApplication(ctx, applicationCode)
	if err != nil {
		return "", err
	}

	provider, err := s.providerFinder.LookupProvider(ctx, providerCode)
	if err != nil {
		return "", err
	}

	sessionID := s.idGenerator.GenerateID()
	oauthCfg := newOauthConfig(provider)

	oauthProvider, err := providers.NewProvider(oauthCfg)
	if err != nil {
		return "", domain.NewAccessDeniedError("failed to create oauth provider: %v", err)
	}

	// save oauthCfg in the session
	us := newUserSession(app.ID.String(), oauthCfg)

	if pErr := s.sessionCache.Put(ctx, sessionID, us, sessionTTL); pErr != nil {
		return "", pErr
	}

	if action == "" {
		action = defaultAction
	}

	state := fmt.Sprintf("%s:%s", action, sessionID)

	// return the auth URL with the session ID embedded in the state parameter
	return oauthProvider.GetAuthURL(ctx, state), nil
}

// Login implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Login(ctx context.Context, code, state string) (string, error) {
	const actionPlusSessionID = 2

	sessionID := state

	if tokens := strings.Split(state, ":"); len(tokens) == actionPlusSessionID {
		sessionID = tokens[1]
	}

	us, err := s.findUserSession(ctx, sessionID)
	if err != nil {
		return "", err
	}

	oauthProvider, err := s.sessionProvider(us)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return "", err
	}

	token, err := oauthProvider.Authorize(ctx, code)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return "", domain.NewAccessDeniedError("failed to authorize: %v", err)
	}

	ui, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return "", domain.NewAccessDeniedError("failed to get user info: %v", err)
	}

	user := toIdentityUser(ui)

	us.updateToken(token)
	us.updateUser(user)

	if pErr := s.sessionCache.Put(ctx, sessionID, us, sessionTTL); pErr != nil {
		return "", pErr
	}

	reqctx.Logger(ctx).Debug().
		Str("sessionId", sessionID).
		Str("applicationId", us.ApplicationID).
		Any("user", us.User).
		Msg("Login completed")

	return sessionID, nil
}

// Session implements the auth.Service interface.
func (s *Service) Session(ctx context.Context, sessionID string) (auth.Session, error) {
	us, err := s.findUserSession(ctx, sessionID)
	if err != nil {
		return auth.Session{}, err
	}

	return auth.Session{
		Header: auth.Header{
			SessionID:     sessionID,
			ApplicationID: us.ApplicationID,
		},
		User: us.User,
	}, nil
}

// Refresh implements the auth.Service interface.
// It returns true if the token was refreshed, false otherwise.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Refresh(ctx context.Context, sessionID string) (bool, error) {
	us, err := s.findUserSession(ctx, sessionID)
	if err != nil {
		return false, err
	}

	oauthProvider, err := s.sessionProvider(us)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return false, err
	}

	token, err := oauthProvider.RefreshAccessToken(ctx, us.Token)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return false, domain.NewAccessDeniedError("failed to refresh token: %v", err)
	}

	if token.AccessToken == us.Token.AccessToken {
		return false, nil
	}

	us.updateToken(token)

	if pErr := s.sessionCache.Put(ctx, sessionID, us, sessionTTL); pErr != nil {
		return false, pErr
	}

	reqctx.Logger(ctx).Debug().
		Str("sessionId", sessionID).
		Msg("Session refreshed")

	return true, nil
}

// Logout implements the auth.Service interface.
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	us, err := s.findUserSession(ctx, sessionID)
	if err != nil {
		return err
	}

	oauthProvider, err := s.sessionProvider(us)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return err
	}

	if rErr := oauthProvider.RevokeAccessToken(ctx, us.Token); rErr != nil {
		return domain.NewAccessDeniedError("failed to revoke token: %v", rErr)
	}

	reqctx.Logger(ctx).Debug().
		Str("sessionId", sessionID).
		Msg("Logout completed")

	return nil
}

// VerifyAPIKey implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) VerifyAPIKey(ctx context.Context, appID admin.ID, apiKey string) error {
	_, err := s.apiKeyFinder.LookupAPIKey(ctx, appID, apiKey)

	return err
}

//nolint:wrapcheck // see comment in the header
func (s *Service) findUserSession(ctx context.Context, sessionID string) (userSession, error) {
	us := userSession{}

	found, err := s.sessionCache.Get(ctx, sessionID, &us)
	if err != nil {
		return userSession{}, err
	}

	if !found {
		return userSession{}, domain.NewAccessDeniedError("invalid userSession ID: %s", sessionID)
	}

	return us, nil
}

//nolint:ireturn
func (s *Service) sessionProvider(session userSession) (oauth.Provider, error) {
	provider, err := providers.NewProvider(session.Config)
	if err != nil {
		return nil, domain.NewAccessDeniedError("failed to create oauth provider: %v", err)
	}

	return provider, nil
}

func (s *Service) silentlyDeleteSession(ctx context.Context, sessionID string) {
	if err := s.sessionCache.Delete(ctx, sessionID); err != nil {
		slog.FromContext(ctx).Info().Err(err).Msg("failed to delete userSession")
	}
}
