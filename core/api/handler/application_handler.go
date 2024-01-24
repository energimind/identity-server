package handler

import (
	"net/http"

	"github.com/energimind/identity-service/core/appl/service"
	"github.com/gin-gonic/gin"
)

// ApplicationHandler is a REST API handler for managing applications.
type ApplicationHandler struct {
	service *service.ApplicationService
}

// NewApplicationHandler creates a new ApplicationHandler.
func NewApplicationHandler(service *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

// Bind binds the ApplicationHandler to a root provided by a router.
func (h *ApplicationHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:id", h.findByID)
	root.POST("", h.create)
	root.PUT("/:id", h.update)
	root.DELETE("/:id", h.delete)
}

func (h *ApplicationHandler) findAll(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (h *ApplicationHandler) findByID(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (h *ApplicationHandler) create(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (h *ApplicationHandler) update(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (h *ApplicationHandler) delete(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
