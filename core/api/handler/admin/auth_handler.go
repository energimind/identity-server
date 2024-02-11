package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// adminActor is the actor for the admin role.
//
//nolint:gochecknoglobals // it is a constant
var adminActor = admin.Actor{Role: admin.SystemRoleAdmin}

// AuthHandler handles admin auth requests.
type AuthHandler struct {
	authEndpoint   string
	userFinder     admin.UserFinder
	cookieOperator admin.CookieOperator
	client         *resty.Client
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(
	authEndpoint string,
	userFinder admin.UserFinder,
	cookieOperator admin.CookieOperator,
) *AuthHandler {
	const clientTimeout = 10 * time.Second

	return &AuthHandler{
		authEndpoint:   authEndpoint + "/api/v1/auth",
		userFinder:     userFinder,
		cookieOperator: cookieOperator,
		client:         resty.New().SetTimeout(clientTimeout),
	}
}

// BindWithMiddlewares binds the LoginHandler to a root provided by a router.
func (h *AuthHandler) BindWithMiddlewares(root gin.IRoutes, mws api.Middlewares) {
	root.GET("/link", h.getProviderLink)
	root.POST("/login", h.completeLogin)
	root.DELETE("/logout", mws.RequireActor, h.logout)
}

func (h *AuthHandler) getProviderLink(c *gin.Context) {
	var result struct {
		Link string `json:"link"`
	}

	appCode := c.Query("appCode")
	providerCode := c.Query("providerCode")

	rsp, err := h.client.R().
		SetContext(c.Request.Context()).
		SetQueryParams(map[string]string{
			"appCode":      appCode,
			"providerCode": providerCode,
		}).
		SetResult(&result).
		Get(h.authEndpoint + "/link")
	if err != nil {
		_ = c.Error(domain.NewGatewayError("failed to get provider link: %v", err))

		return
	}

	if h.processErrorResponse(c, rsp) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"link": result.Link})
}

func (h *AuthHandler) completeLogin(c *gin.Context) { //nolint:funlen
	var completeResult struct {
		SessionID     string `json:"sessionId"`
		ApplicationID string `json:"applicationId"`
		UserInfo      struct {
			Email string `json:"email"`
		} `json:"userInfo"`
	}

	code := c.Query("code")
	state := c.Query("state")

	rsp, err := h.client.R().
		SetContext(c.Request.Context()).
		SetQueryParams(map[string]string{
			"code":  code,
			"state": state,
		}).
		SetResult(&completeResult).
		Post(h.authEndpoint + "/login")
	if err != nil {
		_ = c.Error(domain.NewGatewayError("failed to complete login: %v", err))

		return
	}

	if h.processErrorResponse(c, rsp) {
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

	us := domain.NewUserSession(
		completeResult.SessionID,
		completeResult.ApplicationID,
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
		"applicationId": completeResult.ApplicationID,
		"userId":        user.ID,
		"role":          user.Role,
	})
}

func (h *AuthHandler) logout(c *gin.Context) {
	// reset the cookie even if the logout fails
	if err := h.cookieOperator.ResetCookie(c); err != nil {
		_ = c.Error(err)

		return
	}

	sessionID := c.GetString("sessionId")

	rsp, dErr := h.client.R().
		SetContext(c.Request.Context()).
		SetHeader("X-IS-SessionID", sessionID).
		Delete(h.authEndpoint + "/logout")
	if dErr != nil {
		_ = c.Error(domain.NewGatewayError("failed to logout: %v", dErr))

		return
	}

	if h.processErrorResponse(c, rsp) {
		return
	}

	c.Status(http.StatusOK)
}

func (h *AuthHandler) processErrorResponse(c *gin.Context, rsp *resty.Response) bool {
	if rsp.StatusCode() == http.StatusOK {
		return false
	}

	var result struct {
		Error string `json:"error"`
	}

	// ignore error if the response is not JSON
	_ = json.Unmarshal(rsp.Body(), &result)

	if result.Error == "" {
		result.Error = "unspecified"
	}

	// log this
	_ = c.Error(domain.NewGatewayError("%s (%d)", result.Error, rsp.StatusCode()))

	// send the error to the client, without mapping in the error_mapper middleware
	c.AbortWithStatusJSON(rsp.StatusCode(), gin.H{"error": result.Error})

	return true
}
