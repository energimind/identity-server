package dto

import "time"

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

// User represents an organic user in the system.
type User struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Role        string   `json:"role"`
	APIKeys     []APIKey `json:"apiKeys"`
}

// Daemon represents a non-organic user in the system.
type Daemon struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	APIKeys     []APIKey `json:"apiKeys"`
}

// APIKey represents an API key that can be used to authenticate a daemon.
// It can also be used to authenticate a user.
type APIKey struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Key         string    `json:"key"`
	ExpiresAt   time.Time `json:"expiresAt"`
}
