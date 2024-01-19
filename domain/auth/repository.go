package auth

import "context"

// ApplicationRepository defines the application repository interface.
type ApplicationRepository interface {
	GetApplications(ctx context.Context, principal Principal) ([]Application, error)
	GetApplication(ctx context.Context, principal Principal, id ID) (Application, error)
	CreateApplication(ctx context.Context, principal Principal, app Application) error
	UpdateApplication(ctx context.Context, principal Principal, app Application) error
	DeleteApplication(ctx context.Context, principal Principal, id ID) error
}

// ProviderRepository defines the provider repository interface.
type ProviderRepository interface {
	GetProviders(ctx context.Context, principal Principal, appID string) ([]Provider, error)
	GetProvider(ctx context.Context, principal Principal, appID, id ID) (Provider, error)
	CreateProvider(ctx context.Context, principal Principal, appID string, provider Provider) error
	UpdateProvider(ctx context.Context, principal Principal, appID string, provider Provider) error
	DeleteProvider(ctx context.Context, principal Principal, appID, id ID) error
}

// UserRepository defines the user repository interface.
type UserRepository interface {
	GetUsers(ctx context.Context, principal Principal) ([]User, error)
	GetUser(ctx context.Context, principal Principal, id ID) (User, error)
	CreateUser(ctx context.Context, principal Principal, user User) error
	UpdateUser(ctx context.Context, principal Principal, user User) error
	DeleteUser(ctx context.Context, principal Principal, id ID) error
}

// AccountRepository defines the account repository interface.
type AccountRepository interface {
	GetAccounts(ctx context.Context, principal Principal, userID string) ([]Account, error)
	GetAccount(ctx context.Context, principal Principal, userID, id ID) (Account, error)
	CreateAccount(ctx context.Context, principal Principal, userID string, account Account) error
	UpdateAccount(ctx context.Context, principal Principal, userID string, account Account) error
	DeleteAccount(ctx context.Context, principal Principal, userID, id ID) error
}

// DaemonRepository defines the daemon repository interface.
type DaemonRepository interface {
	GetDaemons(ctx context.Context, principal Principal) ([]Daemon, error)
	GetDaemon(ctx context.Context, principal Principal, id ID) (Daemon, error)
	CreateDaemon(ctx context.Context, principal Principal, daemon Daemon) error
	UpdateDaemon(ctx context.Context, principal Principal, daemon Daemon) error
	DeleteDaemon(ctx context.Context, principal Principal, id ID) error
}

// APIKeyRepository defines the API key repository interface.
type APIKeyRepository interface {
	GetAPIKeys(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID string) ([]APIKey, error)
	GetAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID, id ID) (APIKey, error)
	CreateAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID string, key APIKey) error
	UpdateAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID string, key APIKey) error
	DeleteAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID, id ID) error
}
