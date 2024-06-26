package server

import (
	"context"
	"fmt"

	"github.com/energimind/go-kit/httpd"
	"github.com/energimind/go-kit/idgen/cuuid"
	"github.com/energimind/go-kit/idgen/shortid"
	"github.com/energimind/go-kit/idgen/uuid"
	"github.com/energimind/go-kit/rest/router"
	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/internal/config"
	"github.com/energimind/identity-server/internal/core/api"
	"github.com/energimind/identity-server/internal/core/infra/rest/middleware"
	"github.com/energimind/identity-server/internal/core/infra/rest/sessioncookie"
	"github.com/gin-gonic/gin"
)

// setupServer creates and configures a new server.
// It returns the server, a function to release resources and an error if any.
func setupServer(cfg *config.Config) (*httpd.Server, *closer, error) { //nolint:funlen
	clr := &closer{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startupFailure := func(err error) (*httpd.Server, *closer, error) { //nolint:unparam
		clr.closeAll()

		return nil, nil, err
	}

	idGen := cuuid.NewGenerator()
	shortIDGen := shortid.NewGenerator()
	keyGen := uuid.NewGenerator()

	mongoDB, err := connectMongo(ctx, cfg.Mongo, clr)
	if err != nil {
		return startupFailure(err)
	}

	redisCache, err := connectRedis(ctx, cfg.Redis, clr)
	if err != nil {
		return startupFailure(err)
	}

	cookieOperator := sessioncookie.NewProvider(cfg.Cookie.Name, cfg.Cookie.Secret)

	handlers, middlewares := setupHandlersAndMiddlewares(
		dependencies{
			mongoDB:           mongoDB,
			idGen:             idGen,
			shortIDGen:        shortIDGen,
			keyGen:            keyGen,
			sessionsAPIKey:    cfg.Auth.APIKey,
			localAdminEnabled: cfg.Auth.LocalAdminEnabled,
			cookieOperator:    cookieOperator,
			cache:             redisCache,
		},
	)

	routes := api.NewRoutes(handlers, middlewares)

	restRouter := router.New(
		gin.Recovery(),
		middleware.LoggerInjector(),
		middleware.RequestIDInjector(),
		middleware.RequestLogger(),
		middleware.CORS(cfg.Router.AllowOrigin),
		middleware.ErrorMapper())

	restRouter.RegisterRoutes(routes)

	srv, err := httpd.NewServer(httpd.Config{
		Interface: cfg.HTTP.Interface,
		Port:      cfg.HTTP.Port,
	}, restRouter)
	if err != nil {
		return startupFailure(fmt.Errorf("failed to create server: %w", err))
	}

	slog.Debug().Msgf("Routes:\n%s", formatRoutes(restRouter.GetRoutes()))

	return srv, clr, nil
}
