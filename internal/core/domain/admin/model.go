package admin

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
	SystemRoleManager SystemRole = "manager" // realm management access
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

// Realm represents a realm that can be used to authenticate
// users. It is used to group providers and users.
// It is the top level entity in the admin domain.
type Realm struct {
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
	ID           ID
	Type         ProviderType
	Code         string
	Name         string
	Description  string
	Enabled      bool
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// SystemRole represents the role of a user in the system.
type SystemRole string

// String returns the string representation of the system role.
func (r SystemRole) String() string {
	return string(r)
}

// User represents an organic user in the system.
// The user authenticates with the system using an authentication provider.
type User struct {
	ID          ID
	RealmID     ID
	BindID      string
	Username    string
	Email       string
	DisplayName string
	Description string
	Enabled     bool
	Role        SystemRole
	APIKeys     []APIKey
}

// Daemon represents a non-organic user in the system.
// The daemon authenticates with the system using an API key.
type Daemon struct {
	ID          ID
	RealmID     ID
	Code        string
	Name        string
	Description string
	Enabled     bool
	APIKeys     []APIKey
}

// APIKey represents an API key that can be used to authenticate a daemon.
// It can also be used to authenticate a user.
type APIKey struct {
	ID          ID
	Name        string
	Description string
	Enabled     bool
	Key         string
	ExpiresAt   time.Time
}
