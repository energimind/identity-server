package login

import (
	"context"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/domain/login"
)

// SessionService manages user sessions.
//
// It implements the login.SessionService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type SessionService struct {
	providerLookupService auth.ProviderLookupService
	idgen                 domain.IDGenerator
	sessions              *sessions
}

// NewSessionService returns a new SessionService instance.
func NewSessionService(
	providerLookupService auth.ProviderLookupService,
	idgen domain.IDGenerator,
) *SessionService {
	return &SessionService{
		providerLookupService: providerLookupService,
		idgen:                 idgen,
		sessions:              newSessions(),
	}
}

// Ensure service implements the login.SessionService interface.
var _ login.SessionService = (*SessionService)(nil)

// GetProviderLink implements the login.SessionService interface.
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

	s.sessions.put(sessionID, newSession(cfg))

	return link, nil
}

// CompleteLogin implements the login.SessionService interface.
func (s *SessionService) CompleteLogin(ctx context.Context, code, state string) (string, login.UserInfo, error) {
	sessionID := state

	sess, found := s.sessions.get(sessionID)
	if !found {
		return "", login.UserInfo{}, domain.NewAccessDeniedError("invalid state parameter")
	}

	token, err := exchangeCodeForAccessToken(ctx, sess.Config, code)
	if err != nil {
		s.sessions.delete(sessionID)

		return "", login.UserInfo{}, err
	}

	oui, err := getUserInfo(ctx, token)
	if err != nil {
		s.sessions.delete(sessionID)

		return "", login.UserInfo{}, err
	}

	sess.updateToken(token)

	return sessionID, toUserInfo(oui), nil
}

// Refresh implements the login.SessionService interface.
func (s *SessionService) Refresh(ctx context.Context, sessionID string) error {
	sess, found := s.sessions.get(sessionID)
	if !found {
		return domain.NewAccessDeniedError("invalid session ID")
	}

	token, err := refreshAccessToken(ctx, sess.Config, sess.Token)
	if err != nil {
		return err
	}

	sess.updateToken(token)

	return nil
}

// Logout implements the login.SessionService interface.
func (s *SessionService) Logout(ctx context.Context, sessionID string) error {
	sess, found := s.sessions.get(sessionID)
	if !found {
		return domain.NewAccessDeniedError("invalid session ID")
	}

	s.sessions.delete(sessionID)

	if err := revokeAccessToken(ctx, sess.Token); err != nil {
		return err
	}

	return nil
}
