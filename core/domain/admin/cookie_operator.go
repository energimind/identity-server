package admin

import (
	"github.com/energimind/identity-service/core/domain"
	"github.com/gin-gonic/gin"
)

// CookieOperator defines methods for creating, parsing and resetting cookies.
type CookieOperator interface {
	CreateCookie(c *gin.Context, us domain.UserSession) error
	ResetCookie(c *gin.Context) error
	ParseCookie(c *gin.Context) (domain.UserSession, error)
}
