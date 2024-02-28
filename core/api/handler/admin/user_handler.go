package admin

import (
	"net/http"

	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// UserHandler is an HTTP API handler for managing users.
type UserHandler struct {
	service admin.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service admin.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Bind binds the UserHandler to a root provided by a router.
func (h *UserHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:id", h.findByID)
	root.POST("", h.create)
	root.PUT("/:id", h.update)
	root.DELETE("/:id", h.delete)

	root.GET("/:id/api-keys", h.findAllAPIKeys)
	root.GET("/:id/api-keys/:kid", h.findAPIKey)
	root.POST("/:id/api-keys", h.createAPIKey)
	root.PUT("/:id/api-keys/:kid", h.updateAPIKey)
	root.DELETE("/:id/api-keys/:kid", h.deleteAPIKey)
}

func (h *UserHandler) findAll(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	users, err := h.service.GetUsers(c, actor, admin.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromUsers(users))
}

func (h *UserHandler) findByID(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	user, err := h.service.GetUser(c, actor, admin.ID(appID), admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromUser(user))
}

func (h *UserHandler) create(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoUser := User{}

	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	user := toUser(dtoUser)

	user.ApplicationID = admin.ID(appID)

	user, err := h.service.CreateUser(c, actor, user)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromUser(user))
}

func (h *UserHandler) update(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoUser := User{}

	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	user := toUser(dtoUser)

	user.ID = admin.ID(id)
	user.ApplicationID = admin.ID(appID)

	user, err := h.service.UpdateUser(c, actor, user)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromUser(user))
}

func (h *UserHandler) delete(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteUser(c, actor, admin.ID(appID), admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) findAllAPIKeys(c *gin.Context) {
	appID := c.Param("aid")
	userID := c.Param("id")
	actor := reqctx.Actor(c)

	apiKeys, err := h.service.GetAPIKeys(c, actor, admin.ID(appID), admin.ID(userID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromAPIKeys(apiKeys))
}

func (h *UserHandler) findAPIKey(c *gin.Context) {
	appID := c.Param("aid")
	userID := c.Param("id")
	key := c.Param("kid")
	actor := reqctx.Actor(c)

	apiKey, err := h.service.GetAPIKey(c, actor, admin.ID(appID), admin.ID(userID), admin.ID(key))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromAPIKey(apiKey))
}

func (h *UserHandler) createAPIKey(c *gin.Context) {
	appID := c.Param("aid")
	userID := c.Param("id")
	actor := reqctx.Actor(c)

	dtoAPIKey := APIKey{}

	if err := c.ShouldBindJSON(&dtoAPIKey); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	apiKey := toAPIKey(dtoAPIKey)

	apiKey, err := h.service.CreateAPIKey(c, actor, admin.ID(appID), admin.ID(userID), apiKey)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromAPIKey(apiKey))
}

func (h *UserHandler) updateAPIKey(c *gin.Context) {
	appID := c.Param("aid")
	userID := c.Param("id")
	keyID := c.Param("kid")
	actor := reqctx.Actor(c)

	dtoAPIKey := APIKey{}

	if err := c.ShouldBindJSON(&dtoAPIKey); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	apiKey := toAPIKey(dtoAPIKey)

	apiKey.ID = admin.ID(keyID)

	apiKey, err := h.service.UpdateAPIKey(c, actor, admin.ID(appID), admin.ID(userID), admin.ID(keyID), apiKey)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromAPIKey(apiKey))
}

func (h *UserHandler) deleteAPIKey(c *gin.Context) {
	appID := c.Param("aid")
	userID := c.Param("id")
	keyID := c.Param("kid")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteAPIKey(c, actor, admin.ID(appID), admin.ID(userID), admin.ID(keyID)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
