package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

// ProviderType represents the type of an OAuth provider.
type ProviderType string

// ProviderTypeGoogle represents the Google OAuth provider type.
const ProviderTypeGoogle ProviderType = "google"

// Provider defines the methods that an OAuth provider must implement.
type Provider interface {
	// GetAuthURL returns the URL to redirect the user to in order to authenticate with the provider.
	GetAuthURL(ctx context.Context, state string) string

	// Authorize exchanges the code for an access token.
	Authorize(ctx context.Context, code string) (*oauth2.Token, error)

	// RefreshAccessToken refreshes the access token.
	RefreshAccessToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error)

	// RevokeAccessToken revokes the access token.
	RevokeAccessToken(ctx context.Context, token *oauth2.Token) error

	// GetUserInfo returns the user info.
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}
