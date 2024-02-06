package middleware

import (
	"fmt"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/domain/session"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// RequireActor is a middleware that injects the actor into the request context.
//
// The actor can be retrieved from the request context using the reqctx.Actor function.
//
// If the actor can not be found, the request is aborted with a 401 Unauthorized error.
func RequireActor(verifier auth.CookieVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("sessionKey")
		if err != nil {
			_ = c.Error(domain.NewAccessDeniedError("sessionKey cookie not found"))

			c.Abort()

			return
		}

		serialized, err := verifier.VerifyCookie(cookie)
		if err != nil {
			_ = c.Error(domain.NewAccessDeniedError(fmt.Sprintf("invalid sessionKey cookie: %s", err)))

			c.Abort()

			return
		}

		us, err := session.DeserializeUserSession(serialized)
		if err != nil {
			_ = c.Error(domain.NewAccessDeniedError(fmt.Sprintf("invalid sessionKey cookie value: %s", err)))

			c.Abort()

			return
		}

		actor := auth.NewActor(auth.ID(us.UserID), auth.ID(us.ApplicationID), auth.SystemRole(us.UserRole))

		// add the actor to the request context
		reqctx.SetActor(c, actor)

		// add the actorId to the request context logger
		reqctx.UpdateLogger(c, func(current *zerolog.Logger) zerolog.Logger {
			return current.With().Str("actorId", actor.UserID.String()).Logger()
		})

		c.Next()
	}
}
