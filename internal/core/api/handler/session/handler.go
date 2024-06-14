package session

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/session"
	"github.com/gin-gonic/gin"
)

// Handler is a handler that handles auth requests and sessions.
type Handler struct {
	service    session.Service
	userFinder admin.UserFinder
}

// NewHandler returns a new Handler.
func NewHandler(service session.Service, userFinder admin.UserFinder) *Handler {
	return &Handler{
		service:    service,
		userFinder: userFinder,
	}
}

// Bind binds the Handler to a root provided by a router.
func (h *Handler) Bind(root gin.IRoutes) {
	root.GET("/:sid", h.getSession)
	root.PUT("/:sid/refresh", h.refreshSession)
	root.DELETE("/:sid", h.deleteSession)
	root.GET("/verify", h.verifyAPIKey)
}

// getSession returns the session associated with the session ID.
func (h *Handler) getSession(c *gin.Context) {
	ctx := c.Request.Context()
	sessionID := c.Param("sid")

	sess, err := h.service.Session(ctx, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	realmID := sess.Header.RealmID
	userBindID := sess.User.BindID

	user, err := h.userFinder.GetUserByBindIDSys(ctx, admin.ID(realmID), userBindID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, toClientSession(sess, user))
}

// refreshSession refreshes the session associated with the session ID.
func (h *Handler) refreshSession(c *gin.Context) {
	ctx := c.Request.Context()
	sessionID := c.Param("sid")

	refreshed, err := h.service.Refresh(ctx, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"refreshed": refreshed})
}

// deleteSession logs out the session associated with the session ID,
// and deletes the session from the cache.
func (h *Handler) deleteSession(c *gin.Context) {
	ctx := c.Request.Context()
	sessionID := c.Param("sid")

	err := h.service.Logout(ctx, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}

// verifyAPIKey verifies the API key.
func (h *Handler) verifyAPIKey(c *gin.Context) {
	ctx := c.Request.Context()

	realmID, apiKey, err := decodeAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid authorization header: %v", err))

		return
	}

	err = h.service.VerifyAPIKey(ctx, admin.ID(realmID), apiKey)
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

	const partCount = 2 // realmID:apiKey

	parts := strings.Split(string(decoded), ":")
	if len(parts) != partCount {
		return "", "", fmt.Errorf("invalid API key token format")
	}

	realmID := parts[0]
	apiKey := parts[1]

	return realmID, apiKey, nil
}
