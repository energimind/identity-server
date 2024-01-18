package auth

import "time"

// ProviderTypeGoogle represents the Google authentication provider.
const ProviderTypeGoogle ProviderType = "google"

// SystemRole represents the role of a user in the system.
const (
	SystemRoleUser     SystemRole = "user"     // user only access
	SystemRoleAdmin    SystemRole = "admin"    // application management access
	SystemRoleSysadmin SystemRole = "sysadmin" // system-wide access
)

// KeyOwnerType represents the type of API key owner.
const (
	KeyOwnerTypeUser KeyOwnerType = iota
	KeyOwnerTypeDaemon
)

// Application represents an application that can be used to authenticate
// users. It is used to group providers and users.
// It is the top level entity in the auth domain.
type Application struct {
	ID      ID
	Code    string
	Name    string
	Enabled bool
}

// ProviderType represents the type of authentication provider.
type ProviderType string

// Provider represents an authentication provider.
type Provider struct {
	ID            ID
	ApplicationID ID
	Type          ProviderType
	Code          string
	Name          string
	Enabled       bool
	ClientID      string
	ClientSecret  string
	RedirectURL   string
}

// SystemRole represents the role of a user in the system.
type SystemRole string

// User represents an organic user in the system.
// The user authenticates with the system using an authentication provider.
type User struct {
	ID            ID
	ApplicationID ID
	Username      string
	Enabled       bool
	Role          SystemRole
}

// Account represents an account of a user in the system.
type Account struct {
	ID      ID
	UserID  ID
	UserIDN string
}

// Daemon represents a non-organic user in the system.
// The daemon authenticates with the system using an API key.
type Daemon struct {
	ID            ID
	ApplicationID ID
	Code          string
	Name          string
	Enabled       bool
	Description   string
}

// APIKey represents an API key that can be used to authenticate a daemon.
// It can also be used to authenticate a user.
type APIKey struct {
	ID          ID
	OwnerID     ID
	OwnerType   KeyOwnerType
	Key         string
	Enabled     bool
	Description string
	ExpiresAt   time.Time
}

// KeyOwnerType represents the type of API key owner.
type KeyOwnerType int
