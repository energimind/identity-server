package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is the handler for health check.
//
// It is used to check if the service is ready to serve requests.
type Handler struct{}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Bind binds the Handler to a root provided by a router.
func (h *Handler) Bind(root gin.IRouter) {
	root.GET("/readiness", h.readiness)
	root.GET("/liveness", h.liveness)
}

func (h *Handler) readiness(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *Handler) liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}
