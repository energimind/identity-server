package reqctx

import (
	"github.com/energimind/identity-service/infra/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// SetLogger sets the logger in the given context.
func SetLogger(c *gin.Context, logger zerolog.Logger) {
	// add the logger to the request context
	ctx := logger.WithContext(c.Request.Context())

	// update the request context with the new logger
	c.Request = c.Request.WithContext(ctx)
}

// Logger returns the logger from the given context.
// If the logger was not found in the context, it returns a disabled logger.
func Logger(c *gin.Context) *zerolog.Logger {
	return logger.FromContext(c.Request.Context())
}

// UpdateLogger updates the logger in the given context.
// The update function is called with the current logger.
func UpdateLogger(c *gin.Context, update func(current *zerolog.Logger) zerolog.Logger) {
	// update the logger in the request context
	// propagate the update callback to the logger package
	ctx := logger.UpdateContext(c.Request.Context(), update)

	// update the request context with the new logger
	c.Request = c.Request.WithContext(ctx)
}
