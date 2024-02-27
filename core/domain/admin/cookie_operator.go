package admin

import (
	"github.com/energimind/identity-service/core/domain"
	"github.com/gin-gonic/gin"
)

// CookieOperator defines methods for creating and resetting cookies.
type CookieOperator interface {
	CreateCookie(c *gin.Context, us domain.UserSession) error
	ResetCookie(c *gin.Context) error
}

// CookieParser defines a method for parsing cookies.
type CookieParser interface {
	ParseCookie(c *gin.Context) (domain.UserSession, error)
}
