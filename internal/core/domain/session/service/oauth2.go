package service

import (
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
)

func newOauthConfig(provider admin.Provider) *oauth.Config {
	return &oauth.Config{
		ProviderType: oauth.ProviderType(provider.Type),
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}
