package middleware

import (
	"github.com/energimind/identity-service/core/infra/logger"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// LoggerInjector is a middleware that injects the global logger
// into the request context.
func LoggerInjector() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqctx.SetLogger(c, logger.Global)

		c.Next()
	}
}
