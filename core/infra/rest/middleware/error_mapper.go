package middleware

import (
	"errors"
	"net/http"

	"github.com/energimind/identity-service/core/domain"
	"github.com/gin-gonic/gin"
)

// ErrorMapper is a middleware that maps errors to HTTP responses.
func ErrorMapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// do not map errors if the request has been aborted;
		// assuming the error has already been mapped and the response has been sent
		if c.IsAborted() {
			return
		}

		mapError(c)
	}
}

func mapError(c *gin.Context) {
	var (
		badRequestError   domain.BadRequestError
		accessDeniedError domain.AccessDeniedError
		notFoundError     domain.NotFoundError
		validationError   domain.ValidationError
		storeError        domain.StoreError
		gatewayError      domain.GatewayError
		sessionError      domain.SessionError
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

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
