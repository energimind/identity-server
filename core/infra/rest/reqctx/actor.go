package reqctx

import (
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/gin-gonic/gin"
)

const actorKey = "mw:actor"

// SetActor sets the actor in the given gin context.
func SetActor(c *gin.Context, actor admin.Actor) {
	c.Set(actorKey, actor)
}

// Actor returns the actor from the given gin context.
// The empty actor is returned if the actor was not found in the gin context.
func Actor(c *gin.Context) admin.Actor {
	if value, exists := c.Get(actorKey); exists {
		if actor, ok := value.(admin.Actor); ok {
			return actor
		}
	}

	return admin.Actor{}
}
