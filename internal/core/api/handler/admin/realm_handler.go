package admin

import (
	"net/http"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/infra/rest/reqctx"
	"github.com/gin-gonic/gin"
)

// RealmHandler is a REST API handler for managing realms.
type RealmHandler struct {
	service admin.RealmService
}

// NewRealmHandler creates a new RealmHandler.
func NewRealmHandler(service admin.RealmService) *RealmHandler {
	return &RealmHandler{service: service}
}

// Bind binds the RealmHandler to a root provided by a router.
func (h *RealmHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:aid", h.findByID)
	root.POST("", h.create)
	root.PUT("/:aid", h.update)
	root.DELETE("/:aid", h.delete)
}

func (h *RealmHandler) findAll(c *gin.Context) {
	ctx := c.Request.Context()
	actor := reqctx.Actor(c)

	realms, err := h.service.GetRealms(ctx, actor)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromRealms(realms))
}

func (h *RealmHandler) findByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("aid")
	actor := reqctx.Actor(c)

	realm, err := h.service.GetRealm(ctx, actor, admin.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromRealm(realm))
}

func (h *RealmHandler) create(c *gin.Context) {
	ctx := c.Request.Context()
	actor := reqctx.Actor(c)

	dtoRealm := Realm{}

	if err := c.ShouldBindJSON(&dtoRealm); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid json payload"))

		return
	}

	realm, err := h.service.CreateRealm(ctx, actor, toRealm(dtoRealm))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, fromRealm(realm))
}

func (h *RealmHandler) update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoRealm := Realm{}

	if err := c.ShouldBindJSON(&dtoRealm); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid json payload"))

		return
	}

	dtoRealm.ID = id

	realm, err := h.service.UpdateRealm(ctx, actor, toRealm(dtoRealm))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, fromRealm(realm))
}

func (h *RealmHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("aid")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteRealm(ctx, actor, admin.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
