package login

import (
	"context"
	"time"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/domain/cache"
	"github.com/energimind/identity-service/core/domain/session"
	"github.com/energimind/identity-service/core/infra/logger"
)

const sessionTTL = 24 * 7 * time.Hour

// SessionService manages user sessions.
//
// It implements the session.Service interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type SessionService struct {
	providerLookupService auth.ProviderLookupService
	idgen                 domain.IDGenerator
	cache                 cache.Cache
}

// NewSessionService returns a new SessionService instance.
func NewSessionService(
	providerLookupService auth.ProviderLookupService,
	idgen domain.IDGenerator,
	cache cache.Cache,
) *SessionService {
	return &SessionService{
		providerLookupService: providerLookupService,
		idgen:                 idgen,
		cache:                 cache,
	}
}

// Ensure service implements the session.Service interface.
var _ session.Service = (*SessionService)(nil)

// GetProviderLink implements the session.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *SessionService) GetProviderLink(ctx context.Context, applicationCode, providerCode string) (string, error) {
	provider, err := s.providerLookupService.LookupProvider(ctx, applicationCode, providerCode)
	if err != nil {
		return "", err
	}

	cfg := newConfig(provider)
	sessionID := s.idgen.GenerateID()
	link := getAuthCodeURL(ctx, cfg, sessionID)

	sess := newUserSession(provider.ApplicationID.String(), cfg)

	if pErr := s.cache.Put(ctx, sessionID, sess, sessionTTL); pErr != nil {
		return "", pErr
	}

	return link, nil
}

// CompleteLogin implements the session.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *SessionService) CompleteLogin(ctx context.Context, code, state string) (session.Info, error) {
	sessionID := state

	sess := userSession{}

	found, err := s.cache.Get(ctx, sessionID, &sess)
	if err != nil {
		return session.Info{}, err
	}

	if !found {
		return session.Info{}, domain.NewAccessDeniedError("invalid state parameter")
	}

	token, err := exchangeCodeForAccessToken(ctx, sess.Config, code)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return session.Info{}, err
	}

	oui, err := getUserInfo(ctx, token)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return session.Info{}, err
	}

	sess.updateToken(token)

	if pErr := s.cache.Put(ctx, sessionID, sess, sessionTTL); pErr != nil {
		return session.Info{}, pErr
	}

	info := session.Info{
		SessionID:     sessionID,
		ApplicationID: sess.ApplicationID,
		UserInfo:      toUserInfo(oui),
	}

	return info, nil
}

// Refresh implements the session.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *SessionService) Refresh(ctx context.Context, sessionID string) error {
	sess := userSession{}

	found, err := s.cache.Get(ctx, sessionID, &sess)
	if err != nil {
		return err
	}

	if !found {
		return domain.NewAccessDeniedError("invalid userSession ID")
	}

	token, err := refreshAccessToken(ctx, sess.Config, sess.Token)
	if err != nil {
		return err
	}

	sess.updateToken(token)

	if pErr := s.cache.Put(ctx, sessionID, sess, sessionTTL); pErr != nil {
		return pErr
	}

	return nil
}

// Logout implements the session.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *SessionService) Logout(ctx context.Context, sessionID string) error {
	sess := userSession{}

	found, err := s.cache.Get(ctx, sessionID, &sess)
	if err != nil {
		return err
	}

	if !found {
		return domain.NewAccessDeniedError("invalid userSession ID")
	}

	s.silentlyDeleteSession(ctx, sessionID)

	if rErr := revokeAccessToken(ctx, sess.Token); rErr != nil {
		return rErr
	}

	return nil
}

func (s *SessionService) silentlyDeleteSession(ctx context.Context, sessionID string) {
	if err := s.cache.Delete(ctx, sessionID); err != nil {
		logger.FromContext(ctx).Info().Err(err).Msg("failed to delete userSession")
	}
}
