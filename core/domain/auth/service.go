package auth

import "context"

// Service is a service that handles sessions.
type Service interface {
	ProviderLink(ctx context.Context, applicationCode, providerCode string) (string, error)
	Login(ctx context.Context, code, state string) (Info, error)
	Refresh(ctx context.Context, sessionID string) error
	Logout(ctx context.Context, sessionID string) error
}
