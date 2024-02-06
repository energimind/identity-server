package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const cookieName = "sessionKey"

// UserFinder is an interface for finding users.
type UserFinder interface {
	GetUserByEmail(ctx context.Context, actor admin.Actor, appID admin.ID, email string) (admin.User, error)
}

// AdminLoginHandler handles admin auth requests.
type AdminLoginHandler struct {
	authEndpoint   string
	userFinder     UserFinder
	cookieProvider admin.CookieProvider
	client         *resty.Client
}

// NewAdminLoginHandler returns a new instance of AdminLoginHandler.
func NewAdminLoginHandler(
	authEndpoint string,
	userFinder UserFinder,
	cookieProvider admin.CookieProvider,
) *AdminLoginHandler {
	const clientTimeout = 10 * time.Second

	return &AdminLoginHandler{
		authEndpoint:   authEndpoint + "/api/v1/auth",
		userFinder:     userFinder,
		cookieProvider: cookieProvider,
		client:         resty.New().SetTimeout(clientTimeout),
	}
}

// BindWithMiddlewares binds the LoginHandler to a root provided by a router.
func (h *AdminLoginHandler) BindWithMiddlewares(root gin.IRoutes, mws api.Middlewares) {
	root.GET("/link", h.getProviderLink)
	root.POST("/login", h.completeLogin)
	root.DELETE("/logout", mws.RequireActor, h.logout)
}

func (h *AdminLoginHandler) getProviderLink(c *gin.Context) {
	var result struct {
		Link string `json:"link"`
	}

	appCode := c.Query("appCode")
	providerCode := c.Query("providerCode")

	_, err := h.client.R().
		SetContext(c.Request.Context()).
		SetQueryParams(map[string]string{
			"appCode":      appCode,
			"providerCode": providerCode,
		}).
		SetResult(&result).
		Get(h.authEndpoint + "/link")
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": result.Link})
}

func (h *AdminLoginHandler) completeLogin(c *gin.Context) {
	var completeResult struct {
		SessionID     string `json:"sessionId"`
		ApplicationID string `json:"applicationId"`
		UserInfo      struct {
			Email string `json:"email"`
		} `json:"userInfo"`
	}

	code := c.Query("code")
	state := c.Query("state")

	_, err := h.client.R().
		SetContext(c.Request.Context()).
		SetQueryParams(map[string]string{
			"code":  code,
			"state": state,
		}).
		SetResult(&completeResult).
		Post(h.authEndpoint + "/login")
	if err != nil {
		_ = c.Error(err)

		return
	}

	user, err := h.userFinder.GetUserByEmail(
		c.Request.Context(),
		adminActor,
		admin.ID(completeResult.ApplicationID),
		completeResult.UserInfo.Email,
	)
	if err != nil {
		_ = c.Error(err)

		return
	}

	us := admin.NewUserSession(
		completeResult.SessionID,
		completeResult.ApplicationID,
		user.ID.String(),
		user.Role.String(),
	)

	cookie, err := h.cookieProvider.CreateCookie(c.Request, cookieName, us.Serialize())
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.SetSameSite(cookie.SameSite)
	c.SetCookie(cookieName, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

	c.JSON(http.StatusOK, gin.H{
		"username":    user.Username,
		"email":       user.Email,
		"displayName": user.DisplayName,
		"role":        user.Role,
	})
}

func (h *AdminLoginHandler) logout(c *gin.Context) {
	// reset the cookie even if the logout fails
	cookie, err := h.cookieProvider.ResetCookie(c.Request, cookieName)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.SetSameSite(cookie.SameSite)
	c.SetCookie(cookieName, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

	sessionID := c.GetString("sessionId")

	_, dErr := h.client.R().
		SetContext(c.Request.Context()).
		SetHeader("X-IS-SessionID", sessionID).
		Delete(h.authEndpoint + "/logout")
	if dErr != nil {
		_ = c.Error(dErr)

		return
	}

	c.Status(http.StatusOK)
}
