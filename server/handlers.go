package server

import (
	"github.com/energimind/identity-service/core/api"
	adminapi "github.com/energimind/identity-service/core/api/handler/admin"
	authapi "github.com/energimind/identity-service/core/api/handler/auth"
	healthapi "github.com/energimind/identity-service/core/api/handler/health"
	"github.com/energimind/identity-service/core/appl/service/admin"
	"github.com/energimind/identity-service/core/appl/service/auth"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/infra/cookie"
	"github.com/energimind/identity-service/core/infra/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupHandlers(
	mongoDB *mongo.Database,
	idGen, shortIDGen domain.IDGenerator,
	authEndpoint string,
	cookieProvider *cookie.Provider,
	cache domain.Cache,
) api.Handlers {
	applicationRepo := repository.NewApplicationRepository(mongoDB)
	providerRepo := repository.NewProviderRepository(mongoDB)
	userRepo := repository.NewUserRepository(mongoDB)
	daemonRepo := repository.NewDaemonRepository(mongoDB)

	applicationService := admin.NewApplicationService(applicationRepo, idGen)
	providerService := admin.NewProviderService(providerRepo, idGen)
	userService := admin.NewUserService(userRepo, idGen)
	daemonService := admin.NewDaemonService(daemonRepo, idGen)
	providerLookupService := admin.NewProviderLookupService(applicationService, providerService)
	authService := auth.NewService(providerLookupService, shortIDGen, cache)

	handlers := api.Handlers{
		Application: adminapi.NewApplicationHandler(applicationService),
		Provider:    adminapi.NewProviderHandler(providerService),
		User:        adminapi.NewUserHandler(userService),
		Daemon:      adminapi.NewDaemonHandler(daemonService),
		AdminAuth:   adminapi.NewAuthHandler(authEndpoint, userService, cookieProvider),
		Auth:        authapi.NewHandler(authService),
		Health:      healthapi.NewHandler(),
	}

	return handlers
}
