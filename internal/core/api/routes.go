package api

import (
	"github.com/gin-gonic/gin"
)

// Routes is a collection of routes that will be registered to the router.
type Routes struct {
	handlers    Handlers
	middlewares Middlewares
}

// NewRoutes creates a new Routes.
func NewRoutes(handlers Handlers, middlewares Middlewares) *Routes {
	return &Routes{
		handlers:    handlers,
		middlewares: middlewares,
	}
}

// RegisterRoutes registers the routes to the router.
func (r *Routes) RegisterRoutes(root gin.IRouter) {
	api := root.Group("/api/v1")

	adminEndpoint := api.Group("/admin")
	{
		realmsEndpoint := adminEndpoint.Group("/realms")
		{
			realmsEndpoint.Use(r.middlewares.RequireActor)

			r.bind(realmsEndpoint, r.handlers.Realm)
			r.bind(realmsEndpoint.Group("/:aid/users"), r.handlers.User)
			r.bind(realmsEndpoint.Group("/:aid/daemons"), r.handlers.Daemon)
		}

		providersEndpoint := adminEndpoint.Group("/providers")
		{
			providersEndpoint.Use(r.middlewares.RequireActor)

			r.bind(providersEndpoint, r.handlers.Provider)
		}

		adminAuthEndpoint := adminEndpoint.Group("/auth")
		{
			r.bind(adminAuthEndpoint, r.handlers.Auth)
		}

		utilsEndpoint := adminEndpoint.Group("/utils")
		{
			r.bind(utilsEndpoint, r.handlers.Util)
		}
	}

	sessionsEndpoint := api.Group("/sessions")
	{
		r.bind(sessionsEndpoint, r.handlers.Session)
	}

	healthEndpoint := root.Group("/health")
	{
		r.bind(healthEndpoint, r.handlers.Health)
	}
}

func (r *Routes) bind(root gin.IRoutes, hndlr anyHandler) {
	if mw, ok := hndlr.(handlerWithMiddleware); ok {
		mw.BindWithMiddlewares(root, r.middlewares)

		return
	}

	if h, ok := hndlr.(handler); ok {
		h.Bind(root)

		return
	}

	panic("handler is not a handler or handlerWithMiddleware")
}
