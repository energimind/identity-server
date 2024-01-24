package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler is the handler for health check.
//
// It is used to check if the service is ready to serve requests.
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Bind binds the HealthHandler to a root provided by a router.
func (h *HealthHandler) Bind(root gin.IRoutes) {
	root.GET("/readiness", h.readiness)
	root.GET("/liveness", h.liveness)
}

func (h *HealthHandler) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *HealthHandler) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}
