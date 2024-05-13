package admin

import (
	"net/http"
	"time"

	"github.com/energimind/identity-server/core/api"
	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/domain/local"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	localProviderLink = "/auth/callback?code=" + local.AdminProviderCode +
		"&state=" + local.AdminProviderCode
)

// AuthHandler handles admin auth requests.
type AuthHandler struct {
	identityClient admin.IdentityClient
	cookieOperator admin.CookieOperator
	client         *resty.Client
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(
	identityClient admin.IdentityClient,
	cookieOperator admin.CookieOperator,
) *AuthHandler {
	const clientTimeout = 10 * time.Second

	return &AuthHandler{
		identityClient: identityClient,
		cookieOperator: cookieOperator,
		client:         resty.New().SetTimeout(clientTimeout),
	}
}

// BindWithMiddlewares binds the LoginHandler to a root provided by a router.
func (h *AuthHandler) BindWithMiddlewares(root gin.IRoutes, mws api.Middlewares) {
	root.GET("/link", h.providerLink)
	root.POST("/login", h.login)
	root.DELETE("/logout", mws.RequireActor, h.logout)
}

func (h *AuthHandler) providerLink(c *gin.Context) {
	appCode := c.Query("appCode")
	providerCode := c.Query("providerCode")

	if providerCode == local.AdminProviderCode {
		c.JSON(http.StatusOK, gin.H{"link": localProviderLink})

		return
	}

	link, err := h.identityClient.ProviderLink(c.Request.Context(), appCode, providerCode)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

func (h *AuthHandler) login(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == local.AdminProviderCode && state == local.AdminProviderCode {
		h.loginLocal(c)

		return
	}

	session, user, err := h.identityClient.Login(c.Request.Context(), code, state)
	if err != nil {
		_ = c.Error(err)

		return
	}

	us := domain.NewUserSession(
		session.SessionID,
		session.ApplicationID,
		user.ID.String(),
		user.Role.String(),
	)

	if cErr := h.cookieOperator.CreateCookie(c, us); cErr != nil {
		_ = c.Error(cErr)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":      user.Username,
		"email":         user.Email,
		"displayName":   user.DisplayName,
		"applicationId": session.ApplicationID,
		"userId":        user.ID,
		"role":          user.Role,
	})
}

func (h *AuthHandler) loginLocal(c *gin.Context) {
	us := domain.NewUserSession(
		local.AdminSessionID,
		local.AdminApplicationID,
		local.AdminID,
		local.AdminRole.String(),
	)

	if cErr := h.cookieOperator.CreateCookie(c, us); cErr != nil {
		_ = c.Error(cErr)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":      "admin",
		"email":         "admin",
		"displayName":   "Local Admin",
		"applicationId": local.AdminApplicationID,
		"userId":        local.AdminID,
		"role":          local.AdminRole,
	})
}

func (h *AuthHandler) logout(c *gin.Context) {
	// reset the cookie even if the logout fails
	if err := h.cookieOperator.ResetCookie(c); err != nil {
		_ = c.Error(err)

		return
	}

	sessionID := c.GetString("sessionId")

	if sessionID == local.AdminSessionID {
		c.Status(http.StatusOK)

		return
	}

	if err := h.identityClient.Logout(c.Request.Context(), sessionID); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}
