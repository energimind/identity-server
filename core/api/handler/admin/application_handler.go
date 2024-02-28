package admin

import (
	"net/http"

	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// ApplicationHandler is a REST API handler for managing applications.
type ApplicationHandler struct {
	service admin.ApplicationService
}

// NewApplicationHandler creates a new ApplicationHandler.
func NewApplicationHandler(service admin.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

// Bind binds the ApplicationHandler to a root provided by a router.
func (h *ApplicationHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:aid", h.findByID)
	root.POST("", h.create)
	root.PUT("/:aid", h.update)
	root.DELETE("/:aid", h.delete)
}

func (h *ApplicationHandler) findAll(c *gin.Context) {
	actor := reqctx.Actor(c)

	applications, err := h.service.GetApplications(c, actor)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromApplications(applications))
}

func (h *ApplicationHandler) findByID(c *gin.Context) {
	id := c.Param("aid")
	actor := reqctx.Actor(c)

	application, err := h.service.GetApplication(c, actor, admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromApplication(application))
}

func (h *ApplicationHandler) create(c *gin.Context) {
	actor := reqctx.Actor(c)

	dtoApplication := Application{}

	if err := c.ShouldBindJSON(&dtoApplication); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid json payload"))

		return
	}

	application, err := h.service.CreateApplication(c, actor, toApplication(dtoApplication))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromApplication(application))
}

func (h *ApplicationHandler) update(c *gin.Context) {
	id := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoApplication := Application{}

	if err := c.ShouldBindJSON(&dtoApplication); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid json payload"))

		return
	}

	dtoApplication.ID = id

	application, err := h.service.UpdateApplication(c, actor, toApplication(dtoApplication))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromApplication(application))
}

func (h *ApplicationHandler) delete(c *gin.Context) {
	id := c.Param("aid")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteApplication(c, actor, admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
