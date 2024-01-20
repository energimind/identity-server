package auth

import "context"

// ApplicationRepository defines the application repository interface.
type ApplicationRepository interface {
	GetApplications(ctx context.Context) ([]Application, error)
	GetApplication(ctx context.Context, id ID) (Application, error)
	CreateApplication(ctx context.Context, app Application) error
	UpdateApplication(ctx context.Context, app Application) error
	DeleteApplication(ctx context.Context, id ID) error
}

// ProviderRepository defines the provider repository interface.
type ProviderRepository interface {
	GetProviders(ctx context.Context) ([]Provider, error)
	GetProvider(ctx context.Context, id ID) (Provider, error)
	CreateProvider(ctx context.Context, provider Provider) error
	UpdateProvider(ctx context.Context, provider Provider) error
	DeleteProvider(ctx context.Context, id ID) error
}

// UserRepository defines the user repository interface.
type UserRepository interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id ID) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id ID) error
}

// AccountRepository defines the account repository interface.
type AccountRepository interface {
	GetAccounts(ctx context.Context, userID ID) ([]Account, error)
	GetAccount(ctx context.Context, userID, id ID) (Account, error)
	CreateAccount(ctx context.Context, userID ID, account Account) error
	UpdateAccount(ctx context.Context, userID ID, account Account) error
	DeleteAccount(ctx context.Context, userID, id ID) error
}

// DaemonRepository defines the daemon repository interface.
type DaemonRepository interface {
	GetDaemons(ctx context.Context) ([]Daemon, error)
	GetDaemon(ctx context.Context, id ID) (Daemon, error)
	CreateDaemon(ctx context.Context, daemon Daemon) error
	UpdateDaemon(ctx context.Context, daemon Daemon) error
	DeleteDaemon(ctx context.Context, id ID) error
}

// APIKeyRepository defines the API key repository interface.
type APIKeyRepository interface {
	GetAPIKeys(ctx context.Context, ownerType KeyOwnerType, ownerID ID) ([]APIKey, error)
	GetAPIKey(ctx context.Context, ownerType KeyOwnerType, ownerID, id ID) (APIKey, error)
	CreateAPIKey(ctx context.Context, ownerType KeyOwnerType, ownerID ID, key APIKey) error
	UpdateAPIKey(ctx context.Context, ownerType KeyOwnerType, ownerID ID, key APIKey) error
	DeleteAPIKey(ctx context.Context, ownerType KeyOwnerType, ownerID, id ID) error
}
