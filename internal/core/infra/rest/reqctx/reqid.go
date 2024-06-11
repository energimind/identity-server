package reqctx

import (
	"context"

	"github.com/gin-gonic/gin"
)

// requestIDKey is a context key for the request ID.
type requestIDKey struct{}

// SetRequestID sets the request ID in the underlying request context.
func SetRequestID(c *gin.Context, reqID string) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, requestIDKey{}, reqID)

	c.Request = c.Request.WithContext(ctx)
}

// RequestID returns the request ID from the given context.
// The empty request ID is returned if the request ID was not found in the underlying request context.
//
//nolint:contextcheck // non-inherited new context - not sure why is this a problem
func RequestID(ctx context.Context) string {
	if gctx, ok := ctx.(*gin.Context); ok {
		ctx = gctx.Request.Context()
	}

	if value := ctx.Value(requestIDKey{}); value != nil {
		if reqID, ok := value.(string); ok {
			return reqID
		}
	}

	return ""
}
