package auth

import (
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/infra/oauth"
)

func newOauthConfig(provider admin.Provider) *oauth.Config {
	return &oauth.Config{
		ProviderType: oauth.ProviderType(provider.Type),
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}
