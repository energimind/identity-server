package api

import "github.com/gin-gonic/gin"

// anyHandler represents any handler that can be bound to a router.
type anyHandler interface{}

// handler represents a handler that can be bound to a router.
type handler interface {
	Bind(root gin.IRouter)
}

// handlerWithMiddleware represents a handler that can be bound to a router with middleware.
type handlerWithMiddleware interface {
	BindWithMiddlewares(root gin.IRouter, middlewares Middlewares)
}

// Handlers is a collection of handler that will be bound to the router.
type Handlers struct {
	Auth     anyHandler
	Realm    anyHandler
	Provider anyHandler
	User     anyHandler
	Daemon   anyHandler
	Session  anyHandler
	Util     anyHandler
	Health   anyHandler
}
