package auth

import (
	"context"
	"time"

	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/domain/auth"
	"github.com/energimind/identity-server/core/infra/logger"
	"github.com/energimind/identity-server/core/infra/oauth/providers"
)

const sessionTTL = 24 * time.Hour

// Service manages user sessions.
//
// It implements the auth.Service interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type Service struct {
	providerFinder admin.ProviderLookupService
	apiKeyFinder   admin.APIKeyLookupService
	idGenerator    domain.IDGenerator
	sessionCache   domain.Cache
}

// NewService returns a new Service instance.
func NewService(
	providerFinder admin.ProviderLookupService,
	apiKeyFinder admin.APIKeyLookupService,
	idgen domain.IDGenerator,
	cache domain.Cache,
) *Service {
	return &Service{
		providerFinder: providerFinder,
		apiKeyFinder:   apiKeyFinder,
		idGenerator:    idgen,
		sessionCache:   cache,
	}
}

// Ensure service implements the auth.Service interface.
var _ auth.Service = (*Service)(nil)

// ProviderLink implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) ProviderLink(ctx context.Context, applicationCode, providerCode string) (string, error) {
	provider, err := s.providerFinder.LookupProvider(ctx, applicationCode, providerCode)
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
	session := newUserSession(provider.ApplicationID.String(), oauthCfg)

	if pErr := s.sessionCache.Put(ctx, sessionID, session, sessionTTL); pErr != nil {
		return "", pErr
	}

	// return the auth URL with the session ID embedded in the state parameter
	return oauthProvider.GetAuthURL(ctx, sessionID), nil
}

// Login implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Login(ctx context.Context, code, state string) (auth.Info, error) {
	sessionID := state

	session := userSession{}

	found, err := s.sessionCache.Get(ctx, sessionID, &session)
	if err != nil {
		return auth.Info{}, err
	}

	if !found {
		return auth.Info{}, domain.NewAccessDeniedError("invalid state parameter")
	}

	oauthProvider, err := providers.NewProvider(session.Config)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return auth.Info{}, domain.NewAccessDeniedError("failed to create oauth provider: %v", err)
	}

	token, err := oauthProvider.Authorize(ctx, code)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return auth.Info{}, domain.NewAccessDeniedError("failed to authorize: %v", err)
	}

	ui, err := oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return auth.Info{}, domain.NewAccessDeniedError("failed to get user info: %v", err)
	}

	session.updateToken(token)

	if pErr := s.sessionCache.Put(ctx, sessionID, session, sessionTTL); pErr != nil {
		return auth.Info{}, pErr
	}

	info := auth.Info{
		SessionID:     sessionID,
		ApplicationID: session.ApplicationID,
		UserInfo:      toUserInfo(ui),
	}

	return info, nil
}

// Refresh implements the auth.Service interface.
// It returns true if the token was refreshed, false otherwise.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Refresh(ctx context.Context, sessionID string) (bool, error) {
	session := userSession{}

	found, err := s.sessionCache.Get(ctx, sessionID, &session)
	if err != nil {
		return false, err
	}

	if !found {
		return false, domain.NewAccessDeniedError("invalid userSession ID")
	}

	oauthProvider, err := providers.NewProvider(session.Config)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return false, domain.NewAccessDeniedError("failed to create oauth provider: %v", err)
	}

	token, err := oauthProvider.RefreshAccessToken(ctx, session.Token)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return false, domain.NewAccessDeniedError("failed to refresh token: %v", err)
	}

	if token.AccessToken == session.Token.AccessToken {
		return false, nil
	}

	session.updateToken(token)

	if pErr := s.sessionCache.Put(ctx, sessionID, session, sessionTTL); pErr != nil {
		return false, pErr
	}

	return true, nil
}

// Logout implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	session := userSession{}

	found, err := s.sessionCache.Get(ctx, sessionID, &session)
	if err != nil {
		return err
	}

	if !found {
		return domain.NewAccessDeniedError("invalid userSession ID")
	}

	s.silentlyDeleteSession(ctx, sessionID)

	oauthProvider, err := providers.NewProvider(session.Config)
	if err != nil {
		return domain.NewAccessDeniedError("failed to create oauth provider: %v", err)
	}

	if rErr := oauthProvider.RevokeAccessToken(ctx, session.Token); rErr != nil {
		return domain.NewAccessDeniedError("failed to revoke token: %v", rErr)
	}

	return nil
}

// VerifyAPIKey implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) VerifyAPIKey(ctx context.Context, appID admin.ID, apiKey string) error {
	_, err := s.apiKeyFinder.LookupAPIKey(ctx, appID, apiKey)

	return err
}

func (s *Service) silentlyDeleteSession(ctx context.Context, sessionID string) {
	if err := s.sessionCache.Delete(ctx, sessionID); err != nil {
		logger.FromContext(ctx).Info().Err(err).Msg("failed to delete userSession")
	}
}
