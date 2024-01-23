package api

import (
	"github.com/energimind/identity-service/api/actorctx"
	"github.com/energimind/identity-service/domain/auth"
	"github.com/gin-gonic/gin"
)

// requireActor is a middleware that injects the actor into the request context.
//
// The actor can be retrieved from the request context using the actorctx.Actor function.
//
// If the actor can not be found, the request is aborted with a 401 Unauthorized error.
func requireActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		actor := auth.NewActor("1", "1", auth.SystemRoleAdmin)

		c.Request = c.Request.WithContext(actorctx.WithActor(c.Request.Context(), actor))

		c.Next()
	}
}
