package login

import "github.com/energimind/identity-service/core/domain/auth"

type config struct {
	ProviderType string `json:"providerType"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectUrl"`
}

func newConfig(provider auth.Provider) *config {
	return &config{
		ProviderType: string(provider.Type),
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
	}
}