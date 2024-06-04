package middleware

import (
	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// LoggerInjector is a middleware that injects the global logger
// into the request context.
func LoggerInjector() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqctx.SetLogger(c, slog.Global)

		c.Next()
	}
}
