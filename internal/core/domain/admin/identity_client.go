package admin

import (
	"context"
)

// IdentityClient is an interface that defines the methods to interact with the identity service.
type IdentityClient interface {
	ProviderLink(ctx context.Context, appCode, providerCode string) (string, error)
	Login(ctx context.Context, code, state string) (Session, error)
	Refresh(ctx context.Context, sessionID string) (bool, error)
	Logout(ctx context.Context, sessionID string) error
}

// Session is a struct that contains session information.
type Session struct {
	SessionID     string
	ApplicationID string
	UserEmail     string
}
