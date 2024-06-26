package middleware

import (
	"errors"
	"net/http"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/gin-gonic/gin"
)

// ErrorMapper is a middleware that maps errors to HTTP responses.
func ErrorMapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// do not map errors if the response has been sent already
		if c.Writer.Written() {
			return
		}

		mapError(c)
	}
}

func mapError(c *gin.Context) { //nolint:funlen
	var (
		badRequestError   domain.BadRequestError
		accessDeniedError domain.AccessDeniedError
		notFoundError     domain.NotFoundError
		validationError   domain.ValidationError
		conflictError     domain.ConflictError
		storeError        domain.StoreError
		gatewayError      domain.GatewayError
		sessionError      domain.SessionError
		unauthorizedError domain.UnauthorizedError
	)

	err := c.Errors.Last().Err

	if errors.As(err, &badRequestError) {
		c.JSON(http.StatusBadRequest, gin.H{"error": badRequestError.Error()})

		return
	}

	if errors.As(err, &accessDeniedError) {
		c.JSON(http.StatusForbidden, gin.H{"error": accessDeniedError.Error()})

		return
	}

	if errors.As(err, &notFoundError) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundError.Error()})

		return
	}

	if errors.As(err, &validationError) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})

		return
	}

	if errors.As(err, &conflictError) {
		c.JSON(http.StatusConflict, gin.H{"error": conflictError.Error()})

		return
	}

	if errors.As(err, &storeError) {
		c.JSON(http.StatusBadGateway, gin.H{"error": storeError.Error()})

		return
	}

	if errors.As(err, &gatewayError) {
		c.JSON(http.StatusBadGateway, gin.H{"error": gatewayError.Error()})

		return
	}

	if errors.As(err, &sessionError) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})

		return
	}

	if errors.As(err, &unauthorizedError) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedError.Error()})

		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
