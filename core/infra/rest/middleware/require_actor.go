package middleware

import (
	"fmt"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// RequireActor is a middleware that injects the actor into the request context.
//
// The actor can be retrieved from the request context using the reqctx.Actor function.
//
// If the actor can not be found, the request is aborted with a 401 Unauthorized error.
func RequireActor(verifier admin.CookieVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		us, err := verifier.VerifyCookie(c)
		if err != nil {
			_ = c.Error(domain.NewSessionError(fmt.Sprintf("invalid sessionKey cookie: %s", err)))

			c.Abort()

			return
		}

		// add sessionID to the request
		c.Set("sessionId", us.SessionID)

		actor := admin.NewActor(admin.ID(us.UserID), admin.ID(us.ApplicationID), admin.SystemRole(us.UserRole))

		// add the actor to the request context
		reqctx.SetActor(c, actor)

		// add the actorId to the request context logger
		reqctx.UpdateLogger(c, func(current *zerolog.Logger) zerolog.Logger {
			return current.With().Str("actorId", actor.UserID.String()).Logger()
		})

		c.Next()
	}
}
