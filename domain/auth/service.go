package auth

import "context"

// ApplicationService defines the application service interface.
type ApplicationService interface {
	GetApplications(ctx context.Context, actor Actor) ([]Application, error)
	GetApplication(ctx context.Context, actor Actor, id string) (Application, error)
	CreateApplication(ctx context.Context, actor Actor, app Application) (Application, error)
	UpdateApplication(ctx context.Context, actor Actor, app Application) (Application, error)
	DeleteApplication(ctx context.Context, actor Actor, id string) error
}

// ProviderService defines the provider service interface.
type ProviderService interface {
	GetProviders(ctx context.Context, actor Actor, appID string) ([]Provider, error)
	GetProvider(ctx context.Context, actor Actor, appID, id string) (Provider, error)
	CreateProvider(ctx context.Context, actor Actor, appID string, provider Provider) (Provider, error)
	UpdateProvider(ctx context.Context, actor Actor, appID string, provider Provider) (Provider, error)
	DeleteProvider(ctx context.Context, actor Actor, appID, id string) error
}

// UserService defines the user service interface.
type UserService interface {
	GetUsers(ctx context.Context, actor Actor) ([]User, error)
	GetUser(ctx context.Context, actor Actor, id string) (User, error)
	CreateUser(ctx context.Context, actor Actor, user User) (User, error)
	UpdateUser(ctx context.Context, actor Actor, user User) (User, error)
	DeleteUser(ctx context.Context, actor Actor, id string) error
}

// DaemonService defines the daemon service interface.
type DaemonService interface {
	GetDaemons(ctx context.Context, actor Actor) ([]Daemon, error)
	GetDaemon(ctx context.Context, actor Actor, id string) (Daemon, error)
	CreateDaemon(ctx context.Context, actor Actor, daemon Daemon) (Daemon, error)
	UpdateDaemon(ctx context.Context, actor Actor, daemon Daemon) (Daemon, error)
	DeleteDaemon(ctx context.Context, actor Actor, id string) error
}
