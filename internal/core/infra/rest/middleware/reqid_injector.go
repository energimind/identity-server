package middleware

import (
	"strconv"
	"sync/atomic"

	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// requestIDHeader is the header name for the request ID.
const requestIDHeader = "X-Request-ID"

// RequestIDInjector is a middleware that injects the request ID into the request context.
//
// The request ID can be retrieved from the request context using the reqctx.RequestID function.
func RequestIDInjector() gin.HandlerFunc {
	var idCounter atomic.Int64

	return func(c *gin.Context) {
		// generate a new request ID
		reqID := c.GetHeader(requestIDHeader)

		// reuse the request ID from the header if it is not empty
		if reqID == "" {
			reqID = strconv.Itoa(int(idCounter.Add(1)))
		}

		reqctx.SetRequestID(c, reqID)

		// add the reqId to the request context logger
		reqctx.UpdateLogger(c, func(current *zerolog.Logger) zerolog.Logger {
			return current.With().Str("reqId", reqID).Logger()
		})

		c.Next()
	}
}
