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
		appsEndpoint := adminEndpoint.Group("/applications")
		{
			appsEndpoint.Use(r.middlewares.RequireActor)

			r.bind(appsEndpoint, r.handlers.Application)
			r.bind(appsEndpoint.Group("/:aid/providers"), r.handlers.Provider)
			r.bind(appsEndpoint.Group("/:aid/users"), r.handlers.User)
			r.bind(appsEndpoint.Group("/:aid/daemons"), r.handlers.Daemon)
		}

		adminAuthEndpoint := adminEndpoint.Group("/auth")
		{
			r.bind(adminAuthEndpoint, r.handlers.AdminAuth)
		}

		utilsEndpoint := adminEndpoint.Group("/utils")
		{
			r.bind(utilsEndpoint, r.handlers.Util)
		}
	}

	authEndpoint := api.Group("/auth")
	{
		r.bind(authEndpoint, r.handlers.Auth)
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
