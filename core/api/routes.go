package api

import (
	"github.com/energimind/identity-service/core/infra/rest/middleware"
	"github.com/gin-gonic/gin"
)

// handler represents a handler that can be bound to a router.
type handler interface {
	Bind(root gin.IRoutes)
}

// Handlers is a collection of handler that will be bound to the router.
type Handlers struct {
	Application handler
	Provider    handler
	User        handler
	Daemon      handler
	Auth        handler
	AdminAuth   handler
	Health      handler
}

// Routes is a collection of routes that will be registered to the router.
type Routes struct {
	handlers Handlers
}

// NewRoutes creates a new Routes.
func NewRoutes(handlers Handlers) *Routes {
	return &Routes{
		handlers: handlers,
	}
}

// RegisterRoutes registers the routes to the router.
func (r *Routes) RegisterRoutes(root gin.IRouter) {
	api := root.Group("/api/v1")

	admin := api.Group("/admin")
	{
		admin.Use(middleware.RequireActor())

		apps := admin.Group("/applications")

		r.handlers.Application.Bind(apps)
		r.handlers.Provider.Bind(apps.Group("/:aid/providers"))
		r.handlers.User.Bind(apps.Group("/:aid/users"))
		r.handlers.Daemon.Bind(apps.Group("/:aid/daemons"))

		adminAuth := admin.Group("/auth")

		r.handlers.AdminAuth.Bind(adminAuth)
	}

	auth := api.Group("/auth")
	{
		r.handlers.Auth.Bind(auth)
	}

	health := root.Group("/health")
	{
		r.handlers.Health.Bind(health)
	}
}
