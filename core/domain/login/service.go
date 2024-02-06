package login

import "context"

// SessionService is a service that handles sessions.
type SessionService interface {
	GetProviderLink(ctx context.Context, applicationCode, providerCode string) (string, error)
	CompleteLogin(ctx context.Context, code, state string) (string, UserInfo, error)
	Refresh(ctx context.Context, sessionID string) error
	Logout(ctx context.Context, sessionID string) error
}
