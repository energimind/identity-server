package middleware

import (
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// RequireActor is a middleware that injects the actor into the request context.
//
// The actor can be retrieved from the request context using the reqctx.Actor function.
//
// If the actor can not be found, the request is aborted with a 401 Unauthorized error.
func RequireActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		actor := auth.NewActor("1", "1", auth.SystemRoleAdmin)

		// add the actor to the request context
		reqctx.SetActor(c, actor)

		// add the actorId to the request context logger
		reqctx.UpdateLogger(c, func(current *zerolog.Logger) zerolog.Logger {
			return current.With().Str("actorId", actor.UserID.String()).Logger()
		})

		c.Next()
	}
}
