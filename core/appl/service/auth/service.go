package auth

import (
	"context"
	"time"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/logger"
)

const sessionTTL = 24 * 7 * time.Hour

// Service manages user sessions.
//
// It implements the auth.Service interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type Service struct {
	providerLookupService admin.ProviderLookupService
	idgen                 domain.IDGenerator
	cache                 domain.Cache
}

// NewService returns a new Service instance.
func NewService(
	providerLookupService admin.ProviderLookupService,
	idgen domain.IDGenerator,
	cache domain.Cache,
) *Service {
	return &Service{
		providerLookupService: providerLookupService,
		idgen:                 idgen,
		cache:                 cache,
	}
}

// Ensure service implements the auth.Service interface.
var _ auth.Service = (*Service)(nil)

// GetProviderLink implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) GetProviderLink(ctx context.Context, applicationCode, providerCode string) (string, error) {
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

// CompleteLogin implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) CompleteLogin(ctx context.Context, code, state string) (auth.Info, error) {
	sessionID := state

	sess := userSession{}

	found, err := s.cache.Get(ctx, sessionID, &sess)
	if err != nil {
		return auth.Info{}, err
	}

	if !found {
		return auth.Info{}, domain.NewAccessDeniedError("invalid state parameter")
	}

	token, err := exchangeCodeForAccessToken(ctx, sess.Config, code)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return auth.Info{}, err
	}

	oui, err := getUserInfo(ctx, token)
	if err != nil {
		s.silentlyDeleteSession(ctx, sessionID)

		return auth.Info{}, err
	}

	sess.updateToken(token)

	if pErr := s.cache.Put(ctx, sessionID, sess, sessionTTL); pErr != nil {
		return auth.Info{}, pErr
	}

	info := auth.Info{
		SessionID:     sessionID,
		ApplicationID: sess.ApplicationID,
		UserInfo:      toUserInfo(oui),
	}

	return info, nil
}

// Refresh implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Refresh(ctx context.Context, sessionID string) error {
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

// Logout implements the auth.Service interface.
//
//nolint:wrapcheck // see comment in the header
func (s *Service) Logout(ctx context.Context, sessionID string) error {
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

func (s *Service) silentlyDeleteSession(ctx context.Context, sessionID string) {
	if err := s.cache.Delete(ctx, sessionID); err != nil {
		logger.FromContext(ctx).Info().Err(err).Msg("failed to delete userSession")
	}
}
