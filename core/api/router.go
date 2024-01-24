package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler interface {
	Bind(root gin.IRoutes)
}

// Config is the configuration for the router.
type Config struct {
	AllowOrigin string
}

// Handlers is a collection of handler that will be bound to the router.
type Handlers struct {
	ApplicationHandler handler
	HealthHandler      handler
}

// Router is the router for the REST API.
type Router struct {
	rtr *gin.Engine
}

// NewRouter creates a new Router.
func NewRouter(config Config, handlers Handlers) *Router {
	return &Router{
		rtr: configure(config, handlers),
	}
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}

// GetRoutes returns the routes of the router.
func (r *Router) GetRoutes() []RouteInfo {
	routes := r.rtr.Routes()
	infos := make([]RouteInfo, 0, len(routes))

	for _, route := range routes {
		infos = append(infos, RouteInfo{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	return infos
}

func configure(config Config, handlers Handlers) *gin.Engine {
	rtr := gin.New()

	rtr.Use(loggerInjector())
	rtr.Use(requestLogger())
	rtr.Use(cors(config.AllowOrigin))

	admin := rtr.Group("/api/v1/admin")
	{
		admin.Use(requireActor())

		handlers.ApplicationHandler.Bind(admin.Group("/applications"))
	}

	health := rtr.Group("/health")
	{
		handlers.HealthHandler.Bind(health)
	}

	return rtr
}
