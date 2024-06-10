package reqctx

import (
	"context"

	"github.com/energimind/go-kit/slog"
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

// RequestLogger returns the logger from the given context.
// If the logger was not found in the context, it returns a disabled logger.
func RequestLogger(c *gin.Context) *zerolog.Logger {
	return slog.FromContext(c.Request.Context())
}

// Logger returns the logger from the given context.
// If the logger was not found in the context, it returns a disabled logger.
func Logger(ctx context.Context) *zerolog.Logger {
	return slog.FromContext(ctx)
}

// UpdateLogger updates the logger in the given context.
// The update function is called with the current logger.
func UpdateLogger(c *gin.Context, update func(current *zerolog.Logger) zerolog.Logger) {
	// update the logger in the request context
	// propagate the update callback to the logger package
	ctx := updateContextLogger(c.Request.Context(), update)

	// update the request context with the new logger
	c.Request = c.Request.WithContext(ctx)
}

func updateContextLogger(ctx context.Context, cb func(logger *zerolog.Logger) zerolog.Logger) context.Context {
	existing := slog.FromContext(ctx)
	updated := cb(existing)

	return updated.WithContext(ctx)
}
