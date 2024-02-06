package admin

import (
	"net/http"

	service "github.com/energimind/identity-service/core/appl/service/admin"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// ProviderHandler is an HTTP API handler for managing authentication providers.
type ProviderHandler struct {
	service *service.ProviderService
}

// NewProviderHandler creates a new ProviderHandler.
func NewProviderHandler(service *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: service}
}

// Bind binds the ProviderHandler to a root provided by a router.
func (h *ProviderHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:id", h.findByID)
	root.POST("", h.create)
	root.PUT("/:id", h.update)
	root.DELETE("/:id", h.delete)
}

func (h *ProviderHandler) findAll(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	providers, err := h.service.GetProviders(c, actor, admin.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromProviders(providers))
}

func (h *ProviderHandler) findByID(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	provider, err := h.service.GetProvider(c, actor, admin.ID(appID), admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromProvider(provider))
}

func (h *ProviderHandler) create(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoProvider := Provider{}

	if err := c.ShouldBindJSON(&dtoProvider); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	provider := toProvider(dtoProvider)

	provider.ApplicationID = admin.ID(appID)

	provider, err := h.service.CreateProvider(c, actor, provider)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromProvider(provider))
}

func (h *ProviderHandler) update(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoProvider := Provider{}

	if err := c.ShouldBindJSON(&dtoProvider); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	provider := toProvider(dtoProvider)

	provider.ID = admin.ID(id)
	provider.ApplicationID = admin.ID(appID)

	provider, err := h.service.UpdateProvider(c, actor, provider)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromProvider(provider))
}

func (h *ProviderHandler) delete(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteProvider(c, actor, admin.ID(appID), admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
