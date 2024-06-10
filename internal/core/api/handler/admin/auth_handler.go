package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/api"
	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/local"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	defaultAction     = "login"
	signupAction      = "signup"
	localProviderLink = "/auth/callback?code=" + local.AdminProviderCode + "&state=" + local.AdminProviderCode
)

// adminActor is the actor for the admin role.
//
//nolint:gochecknoglobals // it is a constant
var adminActor = admin.Actor{Role: admin.SystemRoleAdmin}

// AuthHandler handles admin auth requests.
type AuthHandler struct {
	identityClient    *client.Client
	userFinder        admin.UserFinder
	userCreator       admin.UserCreator
	cookieOperator    admin.CookieOperator
	localAdminEnabled bool
	client            *resty.Client
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(
	identityClient *client.Client,
	userFinder admin.UserFinder,
	userCreator admin.UserCreator,
	cookieOperator admin.CookieOperator,
	localAdminEnabled bool,
) *AuthHandler {
	const clientTimeout = 10 * time.Second

	return &AuthHandler{
		identityClient:    identityClient,
		userFinder:        userFinder,
		userCreator:       userCreator,
		cookieOperator:    cookieOperator,
		localAdminEnabled: localAdminEnabled,
		client:            resty.New().SetTimeout(clientTimeout),
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
	action := c.Query("action")

	if action == "" {
		action = defaultAction
	}

	if h.localAdminEnabled && providerCode == local.AdminProviderCode {
		c.JSON(http.StatusOK, gin.H{"link": localProviderLink})

		return
	}

	link, err := h.identityClient.ProviderLink(c.Request.Context(), appCode, providerCode, action)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

func (h *AuthHandler) login(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if strings.HasPrefix(state, signupAction) {
		h.doSignup(c, code, state)

		return
	}

	h.doLogin(c, code, state)
}

func (h *AuthHandler) doLogin(c *gin.Context, code, state string) {
	if h.localAdminEnabled && code == local.AdminProviderCode && state == local.AdminProviderCode {
		h.loginLocal(c)

		return
	}

	session, err := h.identityClient.Login(c.Request.Context(), code, state)
	if err != nil {
		_ = c.Error(err)

		return
	}

	user, err := h.userFinder.GetUserByEmail(
		c.Request.Context(),
		adminActor,
		admin.ID(session.ApplicationID),
		session.User.Email,
	)
	if err != nil {
		_ = c.Error(err)

		return
	}

	h.serveCookieAndUser(c, session, user)
}

func (h *AuthHandler) loginLocal(c *gin.Context) {
	session := client.Session{
		SessionID:     local.AdminSessionID,
		ApplicationID: local.AdminApplicationID,
	}

	user := admin.User{
		ID:            local.AdminID,
		ApplicationID: local.AdminApplicationID,
		Username:      "admin",
		Email:         "admin",
		DisplayName:   "Local Admin",
		Role:          local.AdminRole,
	}

	h.serveCookieAndUser(c, session, user)
}

func (h *AuthHandler) doSignup(c *gin.Context, code, state string) {
	if h.localAdminEnabled && code == local.AdminProviderCode && state == local.AdminProviderCode {
		h.loginLocal(c)

		return
	}

	session, err := h.identityClient.Login(c.Request.Context(), code, state)
	if err != nil {
		_ = c.Error(err)

		return
	}

	oaUser := session.User

	newUser := admin.User{
		ApplicationID: admin.ID(session.ApplicationID),
		Username:      strings.Split(oaUser.Email, "@")[0],
		Email:         oaUser.Email,
		DisplayName:   oaUser.Name,
		Enabled:       true,
		Role:          admin.SystemRoleUser,
	}

	user, err := h.userCreator.CreateUser(c.Request.Context(), adminActor, newUser)
	if err != nil {
		_ = c.Error(err)

		return
	}

	h.serveCookieAndUser(c, session, user)
}

func (h *AuthHandler) serveCookieAndUser(c *gin.Context, session client.Session, user admin.User) {
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

	c.JSON(http.StatusOK, toInfo(session, user))
}

func (h *AuthHandler) logout(c *gin.Context) {
	// reset the cookie even if the logout fails
	if err := h.cookieOperator.ResetCookie(c); err != nil {
		_ = c.Error(err)

		return
	}

	sessionID := c.GetString("sessionId")

	if h.localAdminEnabled && sessionID == local.AdminSessionID {
		c.Status(http.StatusOK)

		return
	}

	if err := h.identityClient.Logout(c.Request.Context(), sessionID); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusOK)
}
