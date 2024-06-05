package server

import (
	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/api"
	adminapi "github.com/energimind/identity-server/internal/core/api/handler/admin"
	authapi "github.com/energimind/identity-server/internal/core/api/handler/auth"
	healthapi "github.com/energimind/identity-server/internal/core/api/handler/health"
	utilapi "github.com/energimind/identity-server/internal/core/api/handler/util"
	"github.com/energimind/identity-server/internal/core/domain"
	adminsvc "github.com/energimind/identity-server/internal/core/domain/admin/service"
	authsvc "github.com/energimind/identity-server/internal/core/domain/auth/service"
	"github.com/energimind/identity-server/internal/core/infra/repository"
	"github.com/energimind/identity-server/internal/core/infra/rest/middleware"
	"github.com/energimind/identity-server/internal/core/infra/rest/sessioncookie"
	"go.mongodb.org/mongo-driver/mongo"
)

type dependencies struct {
	mongoDB           *mongo.Database
	idGen             domain.IDGenerator
	shortIDGen        domain.IDGenerator
	keyGen            domain.IDGenerator
	authEndpoint      string
	localAdminEnabled bool
	cookieOperator    *sessioncookie.Provider
	cache             domain.Cache
}

func setupHandlersAndMiddlewares(deps dependencies) (api.Handlers, api.Middlewares) {
	mongoDB := deps.mongoDB
	idGen := deps.idGen
	shortIDGen := deps.shortIDGen
	keyGen := deps.keyGen
	authEndpoint := deps.authEndpoint
	localAdminEnabled := deps.localAdminEnabled
	cookieOperator := deps.cookieOperator
	cache := deps.cache

	applicationRepo := repository.NewApplicationRepository(mongoDB)
	providerRepo := repository.NewProviderRepository(mongoDB)
	userRepo := repository.NewUserRepository(mongoDB)
	daemonRepo := repository.NewDaemonRepository(mongoDB)

	applicationService := adminsvc.NewApplicationService(applicationRepo, idGen)
	providerService := adminsvc.NewProviderService(providerRepo, idGen)
	userService := adminsvc.NewUserService(userRepo, idGen)
	daemonService := adminsvc.NewDaemonService(daemonRepo, idGen)
	providerLookupService := adminsvc.NewProviderLookupService(applicationService, providerService)
	apiKeyLookupService := adminsvc.NewAPIKeyLookupService(userRepo, daemonRepo)
	authService := authsvc.NewService(providerLookupService, apiKeyLookupService, shortIDGen, cache)

	identityClient := client.New(authEndpoint)

	handlers := api.Handlers{
		Application: adminapi.NewApplicationHandler(applicationService),
		Provider:    adminapi.NewProviderHandler(providerService),
		User:        adminapi.NewUserHandler(userService),
		Daemon:      adminapi.NewDaemonHandler(daemonService),
		AdminAuth:   adminapi.NewAuthHandler(identityClient, userService, userService, cookieOperator, localAdminEnabled),
		Auth:        authapi.NewHandler(authService),
		Util:        utilapi.NewHandler(keyGen),
		Health:      healthapi.NewHandler(),
	}

	middlewares := api.Middlewares{
		RequireActor: middleware.RequireActor(cookieOperator, identityClient, localAdminEnabled),
	}

	return handlers, middlewares
}
