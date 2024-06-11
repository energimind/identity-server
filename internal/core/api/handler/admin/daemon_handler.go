package admin

import (
	"net/http"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// DaemonHandler is an HTTP API handler for managing daemons.
type DaemonHandler struct {
	service admin.DaemonService
}

// NewDaemonHandler creates a new DaemonHandler.
func NewDaemonHandler(service admin.DaemonService) *DaemonHandler {
	return &DaemonHandler{service: service}
}

// Bind binds the DaemonHandler to a root provided by a router.
func (h *DaemonHandler) Bind(root gin.IRoutes) {
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

func (h *DaemonHandler) findAll(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	daemons, err := h.service.GetDaemons(ctx, actor, admin.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromDaemons(daemons))
}

func (h *DaemonHandler) findByID(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	daemon, err := h.service.GetDaemon(ctx, actor, admin.ID(appID), admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromDaemon(daemon))
}

func (h *DaemonHandler) create(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoDaemon := Daemon{}

	if err := c.ShouldBindJSON(&dtoDaemon); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	daemon := toDaemon(dtoDaemon)

	daemon.ApplicationID = admin.ID(appID)

	daemon, err := h.service.CreateDaemon(ctx, actor, daemon)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromDaemon(daemon))
}

func (h *DaemonHandler) update(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoDaemon := Daemon{}

	if err := c.ShouldBindJSON(&dtoDaemon); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	daemon := toDaemon(dtoDaemon)

	daemon.ID = admin.ID(id)
	daemon.ApplicationID = admin.ID(appID)

	daemon, err := h.service.UpdateDaemon(ctx, actor, daemon)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromDaemon(daemon))
}

func (h *DaemonHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteDaemon(ctx, actor, admin.ID(appID), admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DaemonHandler) findAllAPIKeys(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	userID := c.Param("id")
	actor := reqctx.Actor(c)

	apiKeys, err := h.service.GetAPIKeys(ctx, actor, admin.ID(appID), admin.ID(userID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromAPIKeys(apiKeys))
}

func (h *DaemonHandler) findAPIKey(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	userID := c.Param("id")
	keyID := c.Param("kid")
	actor := reqctx.Actor(c)

	apiKey, err := h.service.GetAPIKey(ctx, actor, admin.ID(appID), admin.ID(userID), admin.ID(keyID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromAPIKey(apiKey))
}

func (h *DaemonHandler) createAPIKey(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	userID := c.Param("id")
	actor := reqctx.Actor(c)

	dtoAPIKey := APIKey{}

	if err := c.ShouldBindJSON(&dtoAPIKey); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	apiKey := toAPIKey(dtoAPIKey)

	apiKey, err := h.service.CreateAPIKey(ctx, actor, admin.ID(appID), admin.ID(userID), apiKey)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromAPIKey(apiKey))
}

func (h *DaemonHandler) updateAPIKey(c *gin.Context) {
	ctx := c.Request.Context()
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

	apiKey, err := h.service.UpdateAPIKey(ctx, actor, admin.ID(appID), admin.ID(userID), admin.ID(keyID), apiKey)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromAPIKey(apiKey))
}

func (h *DaemonHandler) deleteAPIKey(c *gin.Context) {
	ctx := c.Request.Context()
	appID := c.Param("aid")
	userID := c.Param("id")
	keyID := c.Param("kid")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteAPIKey(ctx, actor, admin.ID(appID), admin.ID(userID), admin.ID(keyID)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
