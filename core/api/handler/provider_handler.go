package handler

import (
	"net/http"

	"github.com/energimind/identity-service/core/api/dto"
	auth2 "github.com/energimind/identity-service/core/appl/service/auth"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// ProviderHandler is an HTTP API handler for managing authentication providers.
type ProviderHandler struct {
	service *auth2.ProviderService
}

// NewProviderHandler creates a new ProviderHandler.
func NewProviderHandler(service *auth2.ProviderService) *ProviderHandler {
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

	providers, err := h.service.GetProviders(c, actor, auth.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromProviders(providers))
}

func (h *ProviderHandler) findByID(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	provider, err := h.service.GetProvider(c, actor, auth.ID(appID), auth.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromProvider(provider))
}

func (h *ProviderHandler) create(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoProvider := dto.Provider{}

	if err := c.ShouldBindJSON(&dtoProvider); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	provider := dto.ToProvider(dtoProvider)

	provider.ApplicationID = auth.ID(appID)

	provider, err := h.service.CreateProvider(c, actor, provider)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, dto.FromProvider(provider))
}

func (h *ProviderHandler) update(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoProvider := dto.Provider{}

	if err := c.ShouldBindJSON(&dtoProvider); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	provider := dto.ToProvider(dtoProvider)

	provider.ID = auth.ID(id)
	provider.ApplicationID = auth.ID(appID)

	provider, err := h.service.UpdateProvider(c, actor, provider)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromProvider(provider))
}

func (h *ProviderHandler) delete(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteProvider(c, actor, auth.ID(appID), auth.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
