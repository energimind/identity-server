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

	admin := api.Group("/admin")
	{
		apps := admin.Group("/applications")
		{
			apps.Use(r.middlewares.RequireActor)

			r.bind(apps, r.handlers.Application)
			r.bind(apps.Group("/:aid/providers"), r.handlers.Provider)
			r.bind(apps.Group("/:aid/users"), r.handlers.User)
			r.bind(apps.Group("/:aid/daemons"), r.handlers.Daemon)
		}

		adminAuth := admin.Group("/auth")
		{
			r.bind(adminAuth, r.handlers.AdminAuth)
		}
	}

	auth := api.Group("/auth")
	{
		r.bind(auth, r.handlers.Auth)
	}

	health := root.Group("/health")
	{
		r.bind(health, r.handlers.Health)
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
