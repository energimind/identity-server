package admin

import "context"

// RealmRepository defines the realm repository interface.
type RealmRepository interface {
	GetRealms(ctx context.Context) ([]Realm, error)
	GetRealm(ctx context.Context, id ID) (Realm, error)
	CreateRealm(ctx context.Context, realm Realm) error
	UpdateRealm(ctx context.Context, realm Realm) error
	DeleteRealm(ctx context.Context, id ID) error
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
	GetUsers(ctx context.Context, realmID ID) ([]User, error)
	GetUser(ctx context.Context, realmID, id ID) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, realmID, id ID) error
	GetUserByEmail(ctx context.Context, realmID ID, email string) (User, error)
	GetAPIKey(ctx context.Context, realmID ID, key string) (APIKey, error)
}

// DaemonRepository defines the daemon repository interface.
type DaemonRepository interface {
	GetDaemons(ctx context.Context, realmID ID) ([]Daemon, error)
	GetDaemon(ctx context.Context, realmID, id ID) (Daemon, error)
	CreateDaemon(ctx context.Context, daemon Daemon) error
	UpdateDaemon(ctx context.Context, daemon Daemon) error
	DeleteDaemon(ctx context.Context, realmID, id ID) error
	GetAPIKey(ctx context.Context, realmID ID, key string) (APIKey, error)
}
