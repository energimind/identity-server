package api

import (
	"github.com/energimind/identity-service/api/reqctx"
	"github.com/energimind/identity-service/infra/logger"
	"github.com/gin-gonic/gin"
)

// loggerInjector is a middleware that injects the global logger
// into the request context.
func loggerInjector() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqctx.SetLogger(c, logger.Global)

		c.Next()
	}
}
