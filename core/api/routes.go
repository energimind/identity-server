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
	ApplicationHandler handler
	HealthHandler      handler
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
	admin := root.Group("/api/v1/admin")
	{
		admin.Use(middleware.RequireActor())

		r.handlers.ApplicationHandler.Bind(admin.Group("/applications"))
	}

	health := root.Group("/health")
	{
		r.handlers.HealthHandler.Bind(health)
	}
}
