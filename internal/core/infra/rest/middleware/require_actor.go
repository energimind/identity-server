package middleware

import (
	"context"
	"fmt"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/local"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// sessionRefresher is an interface for refreshing a session.
type sessionRefresher interface {
	Refresh(ctx context.Context, sessionID string) (bool, error)
}

// RequireActor is a middleware that injects the actor into the request context.
//
// The actor can be retrieved from the request context using the reqctx.Actor function.
//
// If the actor can not be found, the request is aborted with a 401 Unauthorized error.
//
//nolint:funlen
func RequireActor(
	cookieOperator admin.CookieOperator,
	sessionRefresher sessionRefresher,
	localAdminEnabled bool,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		us, err := cookieOperator.ParseCookie(c)
		if err != nil {
			_ = c.Error(domain.NewSessionError(fmt.Sprintf("invalid sessionKey cookie: %s", err)))

			c.Abort()

			return
		}

		// add sessionID to the request
		c.Set("sessionId", us.SessionID)

		if localAdminEnabled && us.SessionID == local.AdminSessionID && us.UserID == local.AdminID {
			// add the actor to the request context
			reqctx.SetActor(c, admin.NewActor(local.AdminID, local.AdminRealmID, local.AdminRole))

			c.Next()

			return
		}

		refreshed, err := sessionRefresher.Refresh(c, us.SessionID)
		if err != nil {
			// ignoring error, we already have one
			_ = cookieOperator.ResetCookie(c)

			_ = c.Error(domain.NewSessionError("failed to refresh session: %v", err))

			c.Abort()

			return
		}

		if refreshed {
			// update the session cookie
			if err := cookieOperator.CreateCookie(c, us); err != nil {
				_ = c.Error(domain.NewSessionError("failed to update session cookie: %v", err))

				c.Abort()

				return
			}
		}

		actor := admin.NewActor(admin.ID(us.UserID), admin.ID(us.RealmID), admin.SystemRole(us.UserRole))

		// add the actor to the request context
		reqctx.SetActor(c, actor)

		// add the actorId to the request context logger
		reqctx.UpdateLogger(c, func(current *zerolog.Logger) zerolog.Logger {
			return current.With().Str("actorId", actor.UserID.String()).Logger()
		})

		c.Next()
	}
}
