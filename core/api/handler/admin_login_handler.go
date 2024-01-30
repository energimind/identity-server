package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// AdminLoginHandler handles admin login requests.
type AdminLoginHandler struct {
	authEndpoint string
	client       *resty.Client
}

// NewAdminLoginHandler returns a new instance of AdminLoginHandler.
func NewAdminLoginHandler(authEndpoint string) *AdminLoginHandler {
	const clientTimeout = 10 * time.Second

	return &AdminLoginHandler{
		authEndpoint: authEndpoint + "/api/v1/auth",
		client:       resty.New().SetTimeout(clientTimeout),
	}
}

// Bind binds the LoginHandler to a root provided by a router.
func (h *AdminLoginHandler) Bind(root gin.IRoutes) {
	root.GET("/link", h.getProviderLink)
	root.POST("/login", h.completeLogin)
	root.PUT("/refresh", h.refreshSession)
	root.DELETE("/logout", h.logout)
}

func (h *AdminLoginHandler) getProviderLink(c *gin.Context) {
	var result struct {
		Link string `json:"link"`
	}

	_, err := h.client.R().
		SetContext(c.Request.Context()).
		SetQueryParams(map[string]string{
			"appCode":      c.Query("appCode"),
			"providerCode": c.Query("providerCode"),
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
	var result struct {
		SessionID string `json:"sessionId"`
	}

	_, err := h.client.R().
		SetContext(c.Request.Context()).
		SetQueryParams(map[string]string{
			"code":  c.Query("code"),
			"state": c.Query("state"),
		}).
		SetResult(&result).
		Post(h.authEndpoint + "/login")
	if err != nil {
		_ = c.Error(err)

		return
	}

	const cookieMaxAge = 0

	ck := getCookieSecurityContext(c)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(cookieName, encryptCookie(result.SessionID), cookieMaxAge, "/", ck.domain, ck.secure, true)

	c.Status(http.StatusOK)
}

func (h *AdminLoginHandler) refreshSession(c *gin.Context) {
	sessionID := c.GetHeader("X-IS-SessionID")

	_, err := h.client.R().
		SetContext(c.Request.Context()).
		SetHeader("X-IS-SessionID", sessionID).
		Put(h.authEndpoint + "/refresh")
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}

func (h *AdminLoginHandler) logout(c *gin.Context) {
	sessionID := c.GetHeader("X-IS-SessionID")

	_, err := h.client.R().
		SetContext(c.Request.Context()).
		SetHeader("X-IS-SessionID", sessionID).
		Delete(h.authEndpoint + "/logout")
	if err != nil {
		_ = c.Error(err)

		return
	}

	ck := getCookieSecurityContext(c)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(cookieName, "", -1, "/", ck.domain, ck.secure, true)

	c.Status(http.StatusOK)
}
