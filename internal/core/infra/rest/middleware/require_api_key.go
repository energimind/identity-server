package middleware

import (
	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/gin-gonic/gin"
)

// RequireAPIKey is a middleware that requires an API key to be present in the request.
func RequireAPIKey(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		bearer := "Bearer " + apiKey

		if auth != bearer {
			_ = c.Error(domain.NewUnauthorizedError("invalid API key"))

			c.Abort()

			return
		}

		c.Next()
	}
}
