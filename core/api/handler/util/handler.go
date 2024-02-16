package util

import (
	"net/http"

	"github.com/energimind/identity-service/core/domain"
	"github.com/gin-gonic/gin"
)

// Handler is an HTTP handler that provides utility functions.
type Handler struct {
	idgen domain.IDGenerator
}

// NewHandler creates a new Handler.
func NewHandler(idgen domain.IDGenerator) *Handler {
	return &Handler{
		idgen: idgen,
	}
}

// Bind binds the Handler to a root provided by a router.
func (h *Handler) Bind(root gin.IRoutes) {
	root.GET("/key", h.getKey)
}

func (h *Handler) getKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"key": h.idgen.GenerateID()})
}
