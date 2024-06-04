// Package providers provides a builder to create a new
// OAuth provider based on the provider type.
package providers

import (
	"github.com/energimind/identity-server/internal/core/infra/oauth"
	"github.com/energimind/identity-server/internal/core/infra/oauth/providers/google"
)

// NewProvider returns a new OAuth provider based on the provider type.
func NewProvider(config *oauth.Config) (oauth.Provider, error) { //nolint:ireturn
	switch config.ProviderType {
	case oauth.ProviderTypeGoogle:
		return google.NewProvider(config), nil
	default:
		return nil, oauth.NewError("unsupported provider type: %s", config.ProviderType)
	}
}
