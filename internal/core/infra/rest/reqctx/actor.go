package reqctx

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/gin-gonic/gin"
)

// actorKey is a context key for the actor.
type actorKey struct{}

// SetActor sets the actor in the underlying request context.
func SetActor(c *gin.Context, actor admin.Actor) {
	ctx := context.WithValue(c.Request.Context(), actorKey{}, actor)

	c.Request = c.Request.WithContext(ctx)
}

// Actor returns the actor from the given context.
// The empty actor is returned if the actor was not found in the underlying request context.
func Actor(c *gin.Context) admin.Actor {
	if value := c.Request.Context().Value(actorKey{}); value != nil {
		if actor, ok := value.(admin.Actor); ok {
			return actor
		}
	}

	return admin.Actor{}
}
