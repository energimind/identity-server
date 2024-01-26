package dto

// Application represents an application.
type Application struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// Provider represents an authentication provider.
type Provider struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Enabled      bool   `json:"enabled"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectUrl"`
}
