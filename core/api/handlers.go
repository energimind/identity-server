package api

import "github.com/gin-gonic/gin"

// anyHandler represents any handler that can be bound to a router.
type anyHandler interface{}

// handler represents a handler that can be bound to a router.
type handler interface {
	Bind(root gin.IRoutes)
}

// handlerWithMiddleware represents a handler that can be bound to a router with middleware.
type handlerWithMiddleware interface {
	BindWithMiddlewares(root gin.IRoutes, middlewares Middlewares)
}

// Handlers is a collection of handler that will be bound to the router.
type Handlers struct {
	Application anyHandler
	Provider    anyHandler
	User        anyHandler
	Daemon      anyHandler
	AdminAuth   anyHandler
	Auth        anyHandler
	Health      anyHandler
}
