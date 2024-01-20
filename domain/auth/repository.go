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
	GetProviders(ctx context.Context, appID ID) ([]Provider, error)
	GetProvider(ctx context.Context, id ID) (Provider, error)
	CreateProvider(ctx context.Context, provider Provider) error
	UpdateProvider(ctx context.Context, provider Provider) error
	DeleteProvider(ctx context.Context, id ID) error
}

// UserRepository defines the user repository interface.
type UserRepository interface {
	GetUsers(ctx context.Context, appID ID) ([]User, error)
	GetUser(ctx context.Context, id ID) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id ID) error
}

// DaemonRepository defines the daemon repository interface.
type DaemonRepository interface {
	GetDaemons(ctx context.Context, appID ID) ([]Daemon, error)
	GetDaemon(ctx context.Context, id ID) (Daemon, error)
	CreateDaemon(ctx context.Context, daemon Daemon) error
	UpdateDaemon(ctx context.Context, daemon Daemon) error
	DeleteDaemon(ctx context.Context, id ID) error
}
