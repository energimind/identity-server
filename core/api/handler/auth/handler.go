package auth

import (
	"net/http"

	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// Handler is a handler that handles auth requests and sessions.
type Handler struct {
	service auth.Service
}

// NewHandler returns a new Handler.
func NewHandler(service auth.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Bind binds the Handler to a root provided by a router.
func (h *Handler) Bind(root gin.IRoutes) {
	root.GET("/link", h.getProviderLink)
	root.POST("/login", h.completeLogin)
	root.PUT("/refresh", h.refreshSession)
	root.DELETE("/logout", h.logout)
}

func (h *Handler) getProviderLink(c *gin.Context) {
	appCode := c.Query("appCode")
	providerCode := c.Query("providerCode")

	link, err := h.service.ProviderLink(c, appCode, providerCode)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

func (h *Handler) completeLogin(c *gin.Context) {
	info, err := h.service.Login(c, c.Query("code"), c.Query("state"))
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
		"userInfo": gin.H{
			"id":         userInfo.ID,
			"name":       userInfo.Name,
			"given":      userInfo.GivenName,
			"familyName": userInfo.FamilyName,
			"email":      userInfo.Email,
		},
	})
}

func (h *Handler) refreshSession(c *gin.Context) {
	sessionID := c.GetHeader("X-IS-SessionID")

	refreshed, err := h.service.Refresh(c, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"refreshed": refreshed})
}

func (h *Handler) logout(c *gin.Context) {
	sessionID := c.GetHeader("X-IS-SessionID")

	err := h.service.Logout(c, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}
