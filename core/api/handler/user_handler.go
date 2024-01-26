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

// UserHandler is an HTTP API handler for managing users.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Bind binds the UserHandler to a root provided by a router.
func (h *UserHandler) Bind(root gin.IRoutes) {
	root.GET("", h.findAll)
	root.GET("/:id", h.findByID)
	root.POST("", h.create)
	root.PUT("/:id", h.update)
	root.DELETE("/:id", h.delete)
}

func (h *UserHandler) findAll(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	users, err := h.service.GetUsers(c, actor, auth.ID(appID))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromUsers(users))
}

func (h *UserHandler) findByID(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	users, err := h.service.GetUser(c, actor, auth.ID(appID), auth.ID(id))
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromUser(users))
}

func (h *UserHandler) create(c *gin.Context) {
	appID := c.Param("aid")
	actor := reqctx.Actor(c)

	dtoUser := dto.User{}

	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	users := dto.ToUser(dtoUser)

	users.ApplicationID = auth.ID(appID)

	users, err := h.service.CreateUser(c, actor, users)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusCreated, dto.FromUser(users))
}

func (h *UserHandler) update(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	dtoUser := dto.User{}

	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		_ = c.Error(domain.NewBadRequestError("invalid request body"))

		return
	}

	users := dto.ToUser(dtoUser)

	users.ID = auth.ID(id)
	users.ApplicationID = auth.ID(appID)

	users, err := h.service.UpdateUser(c, actor, users)
	if err != nil {
		_ = c.Error(err)

		return
	}

	c.JSON(http.StatusOK, dto.FromUser(users))
}

func (h *UserHandler) delete(c *gin.Context) {
	appID := c.Param("aid")
	id := c.Param("id")
	actor := reqctx.Actor(c)

	if err := h.service.DeleteUser(c, actor, auth.ID(appID), auth.ID(id)); err != nil {
		_ = c.Error(err)

		return
	}

	c.Status(http.StatusNoContent)
}
