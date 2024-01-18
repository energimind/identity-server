package auth

import "context"

// ApplicationRepository defines the application repository interface.
type ApplicationRepository interface {
	GetApplications(ctx context.Context, principal Principal) ([]Application, error)
	GetApplication(ctx context.Context, principal Principal, id string) (Application, error)
	CreateApplication(ctx context.Context, principal Principal, app Application) (Application, error)
	UpdateApplication(ctx context.Context, principal Principal, app Application) (Application, error)
	DeleteApplication(ctx context.Context, principal Principal, id string) error
}

// ProviderRepository defines the provider repository interface.
type ProviderRepository interface {
	GetProviders(ctx context.Context, principal Principal, appID string) ([]Provider, error)
	GetProvider(ctx context.Context, principal Principal, appID, id string) (Provider, error)
	CreateProvider(ctx context.Context, principal Principal, appID string, provider Provider) (Provider, error)
	UpdateProvider(ctx context.Context, principal Principal, appID string, provider Provider) (Provider, error)
	DeleteProvider(ctx context.Context, principal Principal, appID, id string) error
}

// UserRepository defines the user repository interface.
type UserRepository interface {
	GetUsers(ctx context.Context, principal Principal) ([]User, error)
	GetUser(ctx context.Context, principal Principal, id string) (User, error)
	CreateUser(ctx context.Context, principal Principal, user User) (User, error)
	UpdateUser(ctx context.Context, principal Principal, user User) (User, error)
	DeleteUser(ctx context.Context, principal Principal, id string) error
}

// AccountRepository defines the account repository interface.
type AccountRepository interface {
	GetAccounts(ctx context.Context, principal Principal, userID string) ([]Account, error)
	GetAccount(ctx context.Context, principal Principal, userID, id string) (Account, error)
	CreateAccount(ctx context.Context, principal Principal, userID string, account Account) (Account, error)
	UpdateAccount(ctx context.Context, principal Principal, userID string, account Account) (Account, error)
	DeleteAccount(ctx context.Context, principal Principal, userID, id string) error
}

// DaemonRepository defines the daemon repository interface.
type DaemonRepository interface {
	GetDaemons(ctx context.Context, principal Principal) ([]Daemon, error)
	GetDaemon(ctx context.Context, principal Principal, id string) (Daemon, error)
	CreateDaemon(ctx context.Context, principal Principal, daemon Daemon) (Daemon, error)
	UpdateDaemon(ctx context.Context, principal Principal, daemon Daemon) (Daemon, error)
	DeleteDaemon(ctx context.Context, principal Principal, id string) error
}

// APIKeyRepository defines the API key repository interface.
type APIKeyRepository interface {
	GetAPIKeys(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID string) ([]APIKey, error)
	GetAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID, id string) (APIKey, error)
	CreateAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID string, key APIKey) (APIKey, error)
	UpdateAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID string, key APIKey) (APIKey, error)
	DeleteAPIKey(ctx context.Context, principal Principal, ownerType KeyOwnerType, ownerID, id string) error
}
