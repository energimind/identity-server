package auth

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// Service is a service that handles sessions.
type Service interface {
	// Link returns a link to the provider's login page.
	Link(ctx context.Context, realmCode, providerCode, action string) (string, error)

	// Login completes the login/signup process and returns the session ID.
	Login(ctx context.Context, code, state string) (string, error)

	// Session returns the session associated with the session ID.
	Session(ctx context.Context, sessionID string) (Session, error)

	// Refresh refreshes the session associated with the session ID.
	Refresh(ctx context.Context, sessionID string) (bool, error)

	// Logout logs out the session associated with the session ID.
	Logout(ctx context.Context, sessionID string) error

	// VerifyAPIKey verifies the API key.
	VerifyAPIKey(ctx context.Context, realmID admin.ID, apiKey string) error
}
