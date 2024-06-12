package admin

import "context"

// ApplicationService defines the application service interface.
type ApplicationService interface {
	GetApplications(ctx context.Context, actor Actor) ([]Application, error)
	GetApplication(ctx context.Context, actor Actor, id ID) (Application, error)
	CreateApplication(ctx context.Context, actor Actor, app Application) (Application, error)
	UpdateApplication(ctx context.Context, actor Actor, app Application) (Application, error)
	DeleteApplication(ctx context.Context, actor Actor, id ID) error
}

// ProviderService defines the provider service interface.
type ProviderService interface {
	GetProviders(ctx context.Context, actor Actor) ([]Provider, error)
	GetProvider(ctx context.Context, actor Actor, id ID) (Provider, error)
	CreateProvider(ctx context.Context, actor Actor, provider Provider) (Provider, error)
	UpdateProvider(ctx context.Context, actor Actor, provider Provider) (Provider, error)
	DeleteProvider(ctx context.Context, actor Actor, id ID) error
}

// UserService defines the user service interface.
type UserService interface {
	GetUsers(ctx context.Context, actor Actor, appID ID) ([]User, error)
	GetUser(ctx context.Context, actor Actor, appID, id ID) (User, error)
	CreateUser(ctx context.Context, actor Actor, user User) (User, error)
	UpdateUser(ctx context.Context, actor Actor, user User) (User, error)
	DeleteUser(ctx context.Context, actor Actor, appID, id ID) error
	GetAPIKeys(ctx context.Context, actor Actor, appID, userID ID) ([]APIKey, error)
	GetAPIKey(ctx context.Context, actor Actor, appID, userID, id ID) (APIKey, error)
	CreateAPIKey(ctx context.Context, actor Actor, appID, userID ID, apiKey APIKey) (APIKey, error)
	UpdateAPIKey(ctx context.Context, actor Actor, appID, userID, id ID, apiKey APIKey) (APIKey, error)
	DeleteAPIKey(ctx context.Context, actor Actor, appID, userID, id ID) error
}

// UserFinder defines the user finder interface.
// This is a system operation and should not be used in the API.
type UserFinder interface {
	GetUserByEmailSys(ctx context.Context, appID ID, email string) (User, error)
}

// UserCreator defines the user creator interface.
// This is a system operation and should not be used in the API.
type UserCreator interface {
	CreateUserSys(ctx context.Context, user User) (User, error)
}

// DaemonService defines the daemon service interface.
type DaemonService interface {
	GetDaemons(ctx context.Context, actor Actor, appID ID) ([]Daemon, error)
	GetDaemon(ctx context.Context, actor Actor, appID, id ID) (Daemon, error)
	CreateDaemon(ctx context.Context, actor Actor, daemon Daemon) (Daemon, error)
	UpdateDaemon(ctx context.Context, actor Actor, daemon Daemon) (Daemon, error)
	DeleteDaemon(ctx context.Context, actor Actor, appID, id ID) error
	GetAPIKeys(ctx context.Context, actor Actor, appID, daemonID ID) ([]APIKey, error)
	GetAPIKey(ctx context.Context, actor Actor, appID, daemonID, id ID) (APIKey, error)
	CreateAPIKey(ctx context.Context, actor Actor, appID, daemonID ID, apiKey APIKey) (APIKey, error)
	UpdateAPIKey(ctx context.Context, actor Actor, appID, daemonID, id ID, apiKey APIKey) (APIKey, error)
	DeleteAPIKey(ctx context.Context, actor Actor, appID, daemonID, id ID) error
}

// ApplicationLookupService defines the application lookup service interface.
type ApplicationLookupService interface {
	LookupApplication(ctx context.Context, appCode string) (Application, error)
}

// ProviderLookupService defines the provider lookup service interface.
type ProviderLookupService interface {
	LookupProvider(ctx context.Context, providerCode string) (Provider, error)
}

// APIKeyLookupService defines the API key lookup service interface.
type APIKeyLookupService interface {
	LookupAPIKey(ctx context.Context, appID ID, key string) (APIKey, error)
}
