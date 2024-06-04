package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/auth"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
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
	root.GET("/link", h.providerLink)
	root.POST("/login", h.login)
	root.PUT("/refresh", h.refreshSession)
	root.DELETE("/logout", h.logout)
	root.GET("/verify", h.verifyAPIKey)
}

func (h *Handler) providerLink(c *gin.Context) {
	appCode := c.Query("appCode")
	providerCode := c.Query("providerCode")

	link, err := h.service.ProviderLink(c, appCode, providerCode)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

func (h *Handler) login(c *gin.Context) {
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
			"givenName":  userInfo.GivenName,
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

	if refreshed {
		reqctx.Logger(c).Debug().
			Str("sessionId", sessionID).
			Msg("Session refreshed")
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

	reqctx.Logger(c).Debug().
		Str("sessionId", sessionID).
		Msg("Logout completed")

	c.Status(http.StatusOK)
}

func (h *Handler) verifyAPIKey(c *gin.Context) {
	appID, apiKey, err := decodeAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid authorization header: %v", err))

		return
	}

	err = h.service.VerifyAPIKey(c, admin.ID(appID), apiKey)
	if err != nil {
		_ = c.Error(domain.NewUnauthorizedError("invalid API key: %v", err))

		return
	}

	c.Status(http.StatusOK)
}

//nolint:goerr113 // no need to wrap this internal error
func decodeAuthHeader(header string) (string, string, error) {
	if header == "" {
		return "", "", fmt.Errorf("authorization header must not be empty")
	}

	const partCount = 2 // Bearer token

	parts := strings.SplitN(header, " ", partCount)
	if len(parts) != partCount || strings.ToLower(parts[0]) != "bearer" {
		return "", "", fmt.Errorf("invalid authorization header format")
	}

	return decodeAPIKeyToken(parts[1])
}

//nolint:goerr113 // no need to wrap this internal error
func decodeAPIKeyToken(token string) (string, string, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode API key token: %w", err)
	}

	const partCount = 2 // appID:apiKey

	parts := strings.Split(string(decoded), ":")
	if len(parts) != partCount {
		return "", "", fmt.Errorf("invalid API key token format")
	}

	appID := parts[0]
	apiKey := parts[1]

	return appID, apiKey, nil
}
