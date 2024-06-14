package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/energimind/identity-server/internal/core/api"
	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/local"
	"github.com/energimind/identity-server/internal/core/domain/session"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

const (
	defaultAction     = "login"
	signupAction      = "signup"
	localProviderLink = "/auth/callback?code=" + local.AdminProviderCode + "&state=" + local.AdminProviderCode
)

// AuthHandler handles admin auth requests.
type AuthHandler struct {
	sessionService    session.Service
	userFinder        admin.UserFinder
	userCreator       admin.UserCreator
	cookieOperator    admin.CookieOperator
	localAdminEnabled bool
	client            *resty.Client
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(
	sessionService session.Service,
	userFinder admin.UserFinder,
	userCreator admin.UserCreator,
	cookieOperator admin.CookieOperator,
	localAdminEnabled bool,
) *AuthHandler {
	const clientTimeout = 10 * time.Second

	return &AuthHandler{
		sessionService:    sessionService,
		userFinder:        userFinder,
		userCreator:       userCreator,
		cookieOperator:    cookieOperator,
		localAdminEnabled: localAdminEnabled,
		client:            resty.New().SetTimeout(clientTimeout),
	}
}

// BindWithMiddlewares binds the LoginHandler to a root provided by a router.
func (h *AuthHandler) BindWithMiddlewares(root gin.IRoutes, mws api.Middlewares) {
	root.GET("/link", h.link)
	root.POST("/login", h.login)
	root.DELETE("/session", mws.RequireActor, h.logout)
}

// link returns the link to the provider's login page.
func (h *AuthHandler) link(c *gin.Context) {
	realmCode := c.Query("realmCode")
	providerCode := c.Query("providerCode")
	action := c.Query("action")

	if action == "" {
		action = defaultAction
	}

	if h.localAdminEnabled && providerCode == local.AdminProviderCode {
		c.JSON(http.StatusOK, gin.H{"link": localProviderLink})

		return
	}

	ctx := c.Request.Context()

	link, err := h.sessionService.Link(ctx, realmCode, providerCode, action)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

// login completes the login/signup process and returns the session ID.
func (h *AuthHandler) login(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if strings.HasPrefix(state, signupAction) {
		h.doSignup(c, code, state)

		return
	}

	h.doLogin(c, code, state)
}

// logout logs out the user, deletes the session, and resets the cookie.
func (h *AuthHandler) logout(c *gin.Context) {
	if err := h.cookieOperator.ResetCookie(c); err != nil {
		_ = c.Error(err)

		return
	}

	sessionID := c.GetString("sessionId")

	if h.localAdminEnabled && sessionID == local.AdminSessionID {
		c.Status(http.StatusOK)

		return
	}

	ctx := c.Request.Context()

	if err := h.sessionService.Logout(ctx, sessionID); err != nil {
		reqctx.Logger(ctx).Info().
			Str("sessionId", sessionID).
			Err(err).
			Msg("Ignoring failed logout request")
	}

	c.Status(http.StatusOK)
}

func (h *AuthHandler) doLogin(c *gin.Context, code, state string) {
	if h.localAdminEnabled && code == local.AdminProviderCode && state == local.AdminProviderCode {
		h.loginLocal(c)

		return
	}

	ctx := c.Request.Context()

	sessionID, err := h.sessionService.Login(ctx, code, state)
	if err != nil {
		_ = c.Error(err)

		return
	}

	cs, err := h.sessionService.Session(ctx, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	user, err := h.userFinder.GetUserByBindIDSys(ctx, admin.ID(cs.Header.RealmID), cs.User.BindID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	h.serveSessionCookie(c, cs.Header, user)
}

func (h *AuthHandler) loginLocal(c *gin.Context) {
	header := session.Header{
		SessionID: local.AdminSessionID,
		RealmID:   local.AdminRealmID,
	}

	user := admin.User{
		ID:          local.AdminID,
		RealmID:     local.AdminRealmID,
		Username:    "admin",
		Email:       "admin",
		DisplayName: "Local Admin",
		Role:        local.AdminRole,
	}

	h.serveSessionCookie(c, header, user)
}

func (h *AuthHandler) doSignup(c *gin.Context, code, state string) {
	if h.localAdminEnabled && code == local.AdminProviderCode && state == local.AdminProviderCode {
		h.loginLocal(c)

		return
	}

	ctx := c.Request.Context()

	sessionID, err := h.sessionService.Login(ctx, code, state)
	if err != nil {
		_ = c.Error(err)

		return
	}

	cs, err := h.sessionService.Session(ctx, sessionID)
	if err != nil {
		_ = c.Error(err)

		return
	}

	csUser := cs.User

	newUser := admin.User{
		RealmID:     admin.ID(cs.Header.RealmID),
		BindID:      cs.User.BindID,
		Username:    strings.Split(csUser.Email, "@")[0],
		Email:       csUser.Email,
		DisplayName: csUser.Name,
		Enabled:     true,
		Role:        admin.SystemRoleUser,
	}

	user, err := h.userCreator.CreateUserSys(ctx, newUser)
	if err != nil {
		_ = c.Error(err)

		return
	}

	h.serveSessionCookie(c, cs.Header, user)
}

func (h *AuthHandler) serveSessionCookie(c *gin.Context, header session.Header, user admin.User) {
	us := domain.NewUserSession(
		header.SessionID,
		header.RealmID,
		user.ID.String(),
		user.Role.String(),
	)

	if err := h.cookieOperator.CreateCookie(c, us); err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, toSessionInfo(header, user))
}
