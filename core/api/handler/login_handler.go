package handler

import (
	"net/http"

	"github.com/energimind/identity-service/core/domain/login"
	"github.com/gin-gonic/gin"
)

// LoginHandler is a handler that handles login requests and sessions.
type LoginHandler struct {
	service login.SessionService
}

// NewLoginHandler returns a new LoginHandler.
func NewLoginHandler(service login.SessionService) *LoginHandler {
	return &LoginHandler{
		service: service,
	}
}

// Bind binds the LoginHandler to a root provided by a router.
func (h *LoginHandler) Bind(root gin.IRoutes) {
	root.GET("/link", h.getProviderLink)
	root.POST("/login", h.completeLogin)
	root.PUT("/refresh", h.refreshSession)
	root.DELETE("/logout", h.logout)
}

func (h *LoginHandler) getProviderLink(c *gin.Context) {
	appCode := c.Query("appCode")
	providerCode := c.Query("providerCode")

	link, err := h.service.GetProviderLink(c, appCode, providerCode)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

func (h *LoginHandler) completeLogin(c *gin.Context) {
	sessionID, err := h.service.CompleteLogin(c, c.Query("code"), c.Query("state"))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionID})
}

func (h *LoginHandler) refreshSession(c *gin.Context) {
	sessionID := c.GetHeader("X-IS-SessionID")

	err := h.service.Refresh(c, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}

func (h *LoginHandler) logout(c *gin.Context) {
	sessionID := c.GetHeader("X-IS-SessionID")

	err := h.service.Logout(c, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}
