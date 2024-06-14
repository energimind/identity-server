package admin

import (
	"net/http"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// ProviderHandler is an HTTP API handler for managing authentication providers.
type ProviderHandler struct {
	service admin.ProviderService
}

// NewProviderHandler creates a new ProviderHandler.
func NewProviderHandler(service admin.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

// Bind binds the ProviderHandler to a root provided by a router.
func (h *ProviderHandler) Bind(root gin.IRouter) {
	root.GET("", h.findAll)
	root.GET("/:id", h.findByID)
	root.POST("", h.create)
	root.PUT("/:id", h.update)
	root.DELETE("/:id", h.delete)
}

func (h *ProviderHandler) findAll(c *gin.Context) {
	ctx := c.Request.Context()
	actor := reqctx.Actor(c)

	providers, err := h.service.GetProviders(ctx, actor)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromProviders(providers))
}

func (h *ProviderHandler) findByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	actor := reqctx.Actor(c)

	provider, err := h.service.GetProvider(ctx, actor, admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromProvider(provider))
}

func (h *ProviderHandler) create(c *gin.Context) {
	ctx := c.Request.Context()
	actor := reqctx.Actor(c)

	dtoProvider := Provider{}

	if err := c.ShouldBindJSON(&dtoProvider); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	provider := toProvider(dtoProvider)

	provider, err := h.service.CreateProvider(ctx, actor, provider)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromProvider(provider))
}

func (h *ProviderHandler) update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoProvider := Provider{}

	if err := c.ShouldBindJSON(&dtoProvider); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body: %v", err))

		return
	}

	provider := toProvider(dtoProvider)

	provider.ID = admin.ID(id)

	provider, err := h.service.UpdateProvider(ctx, actor, provider)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromProvider(provider))
}

func (h *ProviderHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteProvider(ctx, actor, admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
