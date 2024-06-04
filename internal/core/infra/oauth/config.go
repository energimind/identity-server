package oauth

// Config represents the configuration for an OAuth provider.
type Config struct {
	ProviderType ProviderType `json:"providerType"`
	ClientID     string       `json:"clientId"`
	ClientSecret string       `json:"clientSecret"`
	RedirectURL  string       `json:"redirectUrl"`
}
