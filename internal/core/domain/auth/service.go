package auth

import (
	"context"

	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// Service is a service that handles sessions.
type Service interface {
	ProviderLink(ctx context.Context, applicationCode, providerCode, action string) (string, error)
	Login(ctx context.Context, code, state string) (string, error)
	Session(ctx context.Context, sessionID string) (client.Session, error)
	Refresh(ctx context.Context, sessionID string) (bool, error)
	Logout(ctx context.Context, sessionID string) error
	VerifyAPIKey(ctx context.Context, appID admin.ID, apiKey string) error
}
