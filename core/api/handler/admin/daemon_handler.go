package admin

import (
	"net/http"

	service "github.com/energimind/identity-service/core/appl/service/admin"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// DaemonHandler is an HTTP API handler for managing daemons.
type DaemonHandler struct {
	service *service.DaemonService
}

// NewDaemonHandler creates a new DaemonHandler.
func NewDaemonHandler(service *service.DaemonService) *DaemonHandler {
	return &DaemonHandler{service: service}
}

// Bind binds the DaemonHandler to a root provided by a router.
func (h *DaemonHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:id", h.findByID)
	root.POST("", h.create)
	root.PUT("/:id", h.update)
	root.DELETE("/:id", h.delete)
}

func (h *DaemonHandler) findAll(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	daemonss, err := h.service.GetDaemons(c, actor, admin.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromDaemons(daemonss))
}

func (h *DaemonHandler) findByID(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	daemons, err := h.service.GetDaemon(c, actor, admin.ID(appID), admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromDaemon(daemons))
}

func (h *DaemonHandler) create(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoDaemon := Daemon{}

	if err := c.ShouldBindJSON(&dtoDaemon); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	daemons := toDaemon(dtoDaemon)

	daemons.ApplicationID = admin.ID(appID)

	daemons, err := h.service.CreateDaemon(c, actor, daemons)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromDaemon(daemons))
}

func (h *DaemonHandler) update(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoDaemon := Daemon{}

	if err := c.ShouldBindJSON(&dtoDaemon); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	daemons := toDaemon(dtoDaemon)

	daemons.ID = admin.ID(id)
	daemons.ApplicationID = admin.ID(appID)

	daemons, err := h.service.UpdateDaemon(c, actor, daemons)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromDaemon(daemons))
}

func (h *DaemonHandler) delete(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteDaemon(c, actor, admin.ID(appID), admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
