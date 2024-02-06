package handler

import (
	"net/http"

	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// LoginHandler is a handler that handles auth requests and sessions.
type LoginHandler struct {
	service auth.Service
}

// NewLoginHandler returns a new LoginHandler.
func NewLoginHandler(service auth.Service) *LoginHandler {
	return &LoginHandler{
		service: service,
	}
}

// Bind binds the LoginHandler to a root provided by a router.
func (h *LoginHandler) Bind(root gin.IRoutes) {
	root.GET("/link", h.getProviderLink)
	root.POST("/session", h.completeLogin)
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
	info, err := h.service.CompleteLogin(c, c.Query("code"), c.Query("state"))
	if err != nil {
		_ = c.Error(err)

		return
	}

	userInfo := info.UserInfo

	reqctx.Logger(c).Debug().
		Str("sessionId", info.SessionID).
		Str("applicationId", info.ApplicationID).
		Any("userInfo", userInfo).
		Msg("Login completed")

	c.JSON(http.StatusOK, gin.H{
		"sessionId":     info.SessionID,
		"applicationId": info.ApplicationID,
		"userInfo": map[string]any{
			"id":         userInfo.ID,
			"name":       userInfo.Name,
			"given":      userInfo.GivenName,
			"familyName": userInfo.FamilyName,
			"email":      userInfo.Email,
		},
	})
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
