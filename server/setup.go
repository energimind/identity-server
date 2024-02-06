package server

import (
	"fmt"
	"slices"
	"strings"

	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/api/handler"
	"github.com/energimind/identity-service/core/appl/service/auth"
	"github.com/energimind/identity-service/core/appl/service/login"
	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/cache"
	"github.com/energimind/identity-service/core/infra/cookie"
	"github.com/energimind/identity-service/core/infra/idgen/shortid"
	"github.com/energimind/identity-service/core/infra/idgen/xid"
	"github.com/energimind/identity-service/core/infra/logger"
	"github.com/energimind/identity-service/core/infra/repository"
	"github.com/energimind/identity-service/core/infra/rest/middleware"
	"github.com/energimind/identity-service/core/infra/rest/router"
	"github.com/energimind/identity-service/pkg/httpd"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// setupServer creates and configures a new server.
// It returns the server, a function to release resources and an error if any.
func setupServer(cfg *config.Config) (*httpd.Server, *closer, error) {
	clr := &closer{}

	idGen := xid.NewGenerator()
	shortIDGen := shortid.NewGenerator()

	mongoClient, err := connectToMongoDB(&cfg.Mongo)
	if err != nil {
		return nil, clr, err
	}

	clr.add(func() {
		disconnectFromMongoDB(mongoClient)
	})

	mongoDB := mongoClient.Database(cfg.Mongo.Database)

	redisCache, err := connectToRedis(&cfg.Redis)
	if err != nil {
		return nil, clr, err
	}

	clr.add(func() {
		disconnectFromRedis(redisCache)
	})

	cookieProvider := cookie.NewProvider(cfg.Cookie.Secret)

	handlers := setupHandlers(mongoDB, idGen, shortIDGen, cfg.Auth.Endpoint, cookieProvider, redisCache)

	routes := api.NewRoutes(handlers)

	restRouter := router.New(
		gin.Recovery(),
		middleware.LoggerInjector(),
		middleware.RequestLogger(),
		middleware.CORS(cfg.Router.AllowOrigin),
		middleware.ErrorMapper())

	restRouter.RegisterRoutes(routes)

	srv, err := httpd.NewServer(httpd.Config{
		Interface: cfg.HTTP.Interface,
		Port:      cfg.HTTP.Port,
	}, restRouter)
	if err != nil {
		return nil, clr, fmt.Errorf("failed to create server: %w", err)
	}

	logger.Debug().Msgf("Routes:\n%s", formatRoutes(restRouter.GetRoutes()))

	return srv, clr, nil
}

func setupHandlers(
	mongoDB *mongo.Database,
	idGen, shortIDGen domain.IDGenerator,
	authEndpoint string,
	cookieProvider *cookie.Provider,
	cache cache.Cache,
) api.Handlers {
	applicationRepo := repository.NewApplicationRepository(mongoDB)
	providerRepo := repository.NewProviderRepository(mongoDB)
	userRepo := repository.NewUserRepository(mongoDB)
	daemonRepo := repository.NewDaemonRepository(mongoDB)

	applicationService := auth.NewApplicationService(applicationRepo, idGen)
	providerService := auth.NewProviderService(providerRepo, idGen)
	userService := auth.NewUserService(userRepo, idGen)
	daemonService := auth.NewDaemonService(daemonRepo, idGen)
	providerLookupService := auth.NewProviderLookupService(applicationService, providerService)
	sessionService := login.NewSessionService(providerLookupService, shortIDGen, cache)

	handlers := api.Handlers{
		Application: handler.NewApplicationHandler(applicationService),
		Provider:    handler.NewProviderHandler(providerService),
		User:        handler.NewUserHandler(userService),
		Daemon:      handler.NewDaemonHandler(daemonService),
		Auth:        handler.NewLoginHandler(sessionService),
		AdminAuth:   handler.NewAdminLoginHandler(authEndpoint, userService, cookieProvider),
		Health:      handler.NewHealthHandler(),
	}

	return handlers
}

func formatConfigs(sections []config.Section) string {
	sectionLength := 0

	for _, section := range sections {
		sectionLength = max(sectionLength, len(section.Name))
	}

	var sb strings.Builder

	for i, section := range sections {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(" -> ")
		sb.WriteString(section.Name)
		sb.WriteString(strings.Repeat(" ", sectionLength-len(section.Name)+1))
		sb.WriteRune('{')
		sb.WriteString(strings.Join(section.Values, "; "))
		sb.WriteRune('}')
	}

	return sb.String()
}

func formatRoutes(routes []router.RouteInfo) string {
	methodLength := 0

	for _, route := range routes {
		methodLength = max(methodLength, len(route.Method))
	}

	slices.SortFunc(routes, func(i1, i2 router.RouteInfo) int {
		pd := strings.Compare(i1.Path, i2.Path)

		if pd == 0 {
			// let the order be PUT, GET, DELETE
			return -strings.Compare(i1.Method, i2.Method)
		}

		return pd
	})

	var sb strings.Builder

	for i, route := range routes {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(" -> ")
		sb.WriteString(route.Method)
		sb.WriteString(strings.Repeat(" ", methodLength-len(route.Method)+1))
		sb.WriteString(route.Path)
	}

	return sb.String()
}
