package auth

import "time"

// Provider types.
const (
	ProviderTypeNone   ProviderType = ""
	ProviderTypeGoogle ProviderType = "google"
)

// System roles.
const (
	SystemRoleNone    SystemRole = ""        // no access
	SystemRoleUser    SystemRole = "user"    // user only access
	SystemRoleManager SystemRole = "manager" // application management access
	SystemRoleAdmin   SystemRole = "admin"   // system-wide access
)

// All enums. Used for testing purposes to validate that all enum values are
// covered.
//
//nolint:gochecknoglobals
var (
	AllProviderTypes = []ProviderType{ProviderTypeNone, ProviderTypeGoogle}
	AllSystemRoles   = []SystemRole{SystemRoleNone, SystemRoleUser, SystemRoleManager, SystemRoleAdmin}
)

// Application represents an application that can be used to authenticate
// users. It is used to group providers and users.
// It is the top level entity in the auth domain.
type Application struct {
	ID          ID
	Code        string
	Name        string
	Description string
	Enabled     bool
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
	Description   string
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
	Description   string
	Enabled       bool
	Role          SystemRole
	Accounts      []Account
	APIKeys       []APIKey
}

// Account represents an account of a user in the system.
type Account struct {
	Identifier string
	Enabled    bool
}

// Daemon represents a non-organic user in the system.
// The daemon authenticates with the system using an API key.
type Daemon struct {
	ID            ID
	ApplicationID ID
	Code          string
	Name          string
	Description   string
	Enabled       bool
	APIKeys       []APIKey
}

// APIKey represents an API key that can be used to authenticate a daemon.
// It can also be used to authenticate a user.
type APIKey struct {
	Name        string
	Description string
	Enabled     bool
	Key         string
	ExpiresAt   time.Time
}
