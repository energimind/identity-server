package server

import (
	"github.com/energimind/identity-server/core/api"
	adminapi "github.com/energimind/identity-server/core/api/handler/admin"
	authapi "github.com/energimind/identity-server/core/api/handler/auth"
	healthapi "github.com/energimind/identity-server/core/api/handler/health"
	utilapi "github.com/energimind/identity-server/core/api/handler/util"
	"github.com/energimind/identity-server/core/appl/service/admin"
	"github.com/energimind/identity-server/core/appl/service/auth"
	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/infra/identity"
	"github.com/energimind/identity-server/core/infra/repository"
	"github.com/energimind/identity-server/core/infra/rest/middleware"
	"github.com/energimind/identity-server/core/infra/rest/sessioncookie"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupHandlersAndMiddlewares(
	mongoDB *mongo.Database,
	idGen, shortIDGen, keyGen domain.IDGenerator,
	authEndpoint string,
	cookieOperator *sessioncookie.Provider,
	cache domain.Cache,
) (api.Handlers, api.Middlewares) {
	applicationRepo := repository.NewApplicationRepository(mongoDB)
	providerRepo := repository.NewProviderRepository(mongoDB)
	userRepo := repository.NewUserRepository(mongoDB)
	daemonRepo := repository.NewDaemonRepository(mongoDB)

	applicationService := admin.NewApplicationService(applicationRepo, idGen)
	providerService := admin.NewProviderService(providerRepo, idGen)
	userService := admin.NewUserService(userRepo, idGen)
	daemonService := admin.NewDaemonService(daemonRepo, idGen)
	providerLookupService := admin.NewProviderLookupService(applicationService, providerService)
	apiKeyLookupService := admin.NewAPIKeyLookupService(userRepo, daemonRepo)
	authService := auth.NewService(providerLookupService, apiKeyLookupService, shortIDGen, cache)

	identityClient := identity.NewClient(authEndpoint, userService)

	handlers := api.Handlers{
		Application: adminapi.NewApplicationHandler(applicationService),
		Provider:    adminapi.NewProviderHandler(providerService),
		User:        adminapi.NewUserHandler(userService),
		Daemon:      adminapi.NewDaemonHandler(daemonService),
		AdminAuth:   adminapi.NewAuthHandler(identityClient, cookieOperator),
		Auth:        authapi.NewHandler(authService),
		Util:        utilapi.NewHandler(keyGen),
		Health:      healthapi.NewHandler(),
	}

	middlewares := api.Middlewares{
		RequireActor: middleware.RequireActor(cookieOperator, identityClient),
	}

	return handlers, middlewares
}
