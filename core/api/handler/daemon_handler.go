package handler

import (
	"net/http"

	"github.com/energimind/identity-service/core/api/dto"
	"github.com/energimind/identity-service/core/appl/service"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
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

	daemonss, err := h.service.GetDaemons(c, actor, auth.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromDaemons(daemonss))
}

func (h *DaemonHandler) findByID(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	daemons, err := h.service.GetDaemon(c, actor, auth.ID(appID), auth.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromDaemon(daemons))
}

func (h *DaemonHandler) create(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoDaemon := dto.Daemon{}

	if err := c.ShouldBindJSON(&dtoDaemon); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	daemons := dto.ToDaemon(dtoDaemon)

	daemons.ApplicationID = auth.ID(appID)

	daemons, err := h.service.CreateDaemon(c, actor, daemons)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, dto.FromDaemon(daemons))
}

func (h *DaemonHandler) update(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoDaemon := dto.Daemon{}

	if err := c.ShouldBindJSON(&dtoDaemon); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	daemons := dto.ToDaemon(dtoDaemon)

	daemons.ID = auth.ID(id)
	daemons.ApplicationID = auth.ID(appID)

	daemons, err := h.service.UpdateDaemon(c, actor, daemons)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromDaemon(daemons))
}

func (h *DaemonHandler) delete(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteDaemon(c, actor, auth.ID(appID), auth.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
