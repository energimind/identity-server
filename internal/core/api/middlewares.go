package api

import "github.com/gin-gonic/gin"

// Middlewares is a collection of middlewares that will be bound to the router.
type Middlewares struct {
	RequireActor  gin.HandlerFunc
	RequireAPIKey gin.HandlerFunc
}
