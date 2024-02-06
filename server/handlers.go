package server

import (
	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/api/handler"
	"github.com/energimind/identity-service/core/appl/service/auth"
	"github.com/energimind/identity-service/core/appl/service/login"
	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/cache"
	"github.com/energimind/identity-service/core/infra/cookie"
	"github.com/energimind/identity-service/core/infra/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

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
