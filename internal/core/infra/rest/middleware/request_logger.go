package middleware

import (
	"strings"
	"time"

	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// RequestLogger is a middleware that logs requests.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()

		if shouldNotLogRequest(c) {
			return
		}

		emitLogEntry(c, end.Sub(start))
	}
}

func emitLogEntry(c *gin.Context, duration time.Duration) {
	status := c.Writer.Status()

	event := reqctx.RequestLogger(c).Debug().
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Dur("duration", duration).
		Int("status", status)

	if err := c.Errors.Last(); err != nil {
		event.Err(c.Errors.Last())
	}

	if isFailureStatusCode(status) {
		event.Msg("Request failed")

		return
	}

	event.Msg("Request completed")
}

func shouldNotLogRequest(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/health/")
}

func isFailureStatusCode(statusCode int) bool {
	const failureStatusCodeThreshold = 400

	return statusCode >= failureStatusCodeThreshold
}
