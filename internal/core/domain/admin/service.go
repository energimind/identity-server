package admin

import "context"

// RealmService defines the realm service interface.
type RealmService interface {
	GetRealms(ctx context.Context, actor Actor) ([]Realm, error)
	GetRealm(ctx context.Context, actor Actor, id ID) (Realm, error)
	CreateRealm(ctx context.Context, actor Actor, realm Realm) (Realm, error)
	UpdateRealm(ctx context.Context, actor Actor, realm Realm) (Realm, error)
	DeleteRealm(ctx context.Context, actor Actor, id ID) error
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
	GetUsers(ctx context.Context, actor Actor, realmID ID) ([]User, error)
	GetUser(ctx context.Context, actor Actor, realmID, id ID) (User, error)
	CreateUser(ctx context.Context, actor Actor, user User) (User, error)
	UpdateUser(ctx context.Context, actor Actor, user User) (User, error)
	DeleteUser(ctx context.Context, actor Actor, realmID, id ID) error
	GetAPIKeys(ctx context.Context, actor Actor, realmID, userID ID) ([]APIKey, error)
	GetAPIKey(ctx context.Context, actor Actor, realmID, userID, id ID) (APIKey, error)
	CreateAPIKey(ctx context.Context, actor Actor, realmID, userID ID, apiKey APIKey) (APIKey, error)
	UpdateAPIKey(ctx context.Context, actor Actor, realmID, userID, id ID, apiKey APIKey) (APIKey, error)
	DeleteAPIKey(ctx context.Context, actor Actor, realmID, userID, id ID) error
}

// UserFinder defines the user finder interface.
// This is a system operation and should not be used in the API.
type UserFinder interface {
	GetUserByBindIDSys(ctx context.Context, realmID ID, bindID string) (User, error)
}

// UserCreator defines the user creator interface.
// This is a system operation and should not be used in the API.
type UserCreator interface {
	CreateUserSys(ctx context.Context, user User) (User, error)
}

// DaemonService defines the daemon service interface.
type DaemonService interface {
	GetDaemons(ctx context.Context, actor Actor, realmID ID) ([]Daemon, error)
	GetDaemon(ctx context.Context, actor Actor, realmID, id ID) (Daemon, error)
	CreateDaemon(ctx context.Context, actor Actor, daemon Daemon) (Daemon, error)
	UpdateDaemon(ctx context.Context, actor Actor, daemon Daemon) (Daemon, error)
	DeleteDaemon(ctx context.Context, actor Actor, realmID, id ID) error
	GetAPIKeys(ctx context.Context, actor Actor, realmID, daemonID ID) ([]APIKey, error)
	GetAPIKey(ctx context.Context, actor Actor, realmID, daemonID, id ID) (APIKey, error)
	CreateAPIKey(ctx context.Context, actor Actor, realmID, daemonID ID, apiKey APIKey) (APIKey, error)
	UpdateAPIKey(ctx context.Context, actor Actor, realmID, daemonID, id ID, apiKey APIKey) (APIKey, error)
	DeleteAPIKey(ctx context.Context, actor Actor, realmID, daemonID, id ID) error
}

// RealmLookupService defines the realm lookup service interface.
type RealmLookupService interface {
	LookupRealm(ctx context.Context, realmCode string) (Realm, error)
}

// ProviderLookupService defines the provider lookup service interface.
type ProviderLookupService interface {
	LookupProvider(ctx context.Context, providerCode string) (Provider, error)
}

// APIKeyLookupService defines the API key lookup service interface.
type APIKeyLookupService interface {
	LookupAPIKey(ctx context.Context, realmID ID, key string) (APIKey, error)
}
