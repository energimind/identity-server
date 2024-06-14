package server

import (
	"github.com/energimind/identity-server/internal/core/api"
	adminapi "github.com/energimind/identity-server/internal/core/api/handler/admin"
	healthapi "github.com/energimind/identity-server/internal/core/api/handler/health"
	sessionapi "github.com/energimind/identity-server/internal/core/api/handler/session"
	utilapi "github.com/energimind/identity-server/internal/core/api/handler/util"
	"github.com/energimind/identity-server/internal/core/domain"
	adminsvc "github.com/energimind/identity-server/internal/core/domain/admin/service"
	authsvc "github.com/energimind/identity-server/internal/core/domain/session/service"
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
	sessionsAPIKey    string
	localAdminEnabled bool
	cookieOperator    *sessioncookie.Provider
	cache             domain.Cache
}

func setupHandlersAndMiddlewares(deps dependencies) (api.Handlers, api.Middlewares) {
	mongoDB := deps.mongoDB
	idGen := deps.idGen
	shortIDGen := deps.shortIDGen
	keyGen := deps.keyGen
	sessionAPIKey := deps.sessionsAPIKey
	localAdminEnabled := deps.localAdminEnabled
	cookieOperator := deps.cookieOperator
	cache := deps.cache

	realmRepo := repository.NewRealmRepository(mongoDB)
	providerRepo := repository.NewProviderRepository(mongoDB)
	userRepo := repository.NewUserRepository(mongoDB)
	daemonRepo := repository.NewDaemonRepository(mongoDB)

	realmService := adminsvc.NewRealmService(realmRepo, idGen)
	providerService := adminsvc.NewProviderService(providerRepo, idGen)
	userService := adminsvc.NewUserService(userRepo, idGen)
	daemonService := adminsvc.NewDaemonService(daemonRepo, idGen)
	realmLookupService := adminsvc.NewRealmLookupService(realmService)
	providerLookupService := adminsvc.NewProviderLookupService(providerService)
	apiKeyLookupService := adminsvc.NewAPIKeyLookupService(userRepo, daemonRepo)
	sessionService := authsvc.NewService(
		realmLookupService,
		providerLookupService,
		apiKeyLookupService,
		shortIDGen,
		cache,
	)

	handlers := api.Handlers{
		Auth:     adminapi.NewAuthHandler(sessionService, userService, cookieOperator, localAdminEnabled),
		Realm:    adminapi.NewRealmHandler(realmService),
		Provider: adminapi.NewProviderHandler(providerService),
		User:     adminapi.NewUserHandler(userService),
		Daemon:   adminapi.NewDaemonHandler(daemonService),
		Session:  sessionapi.NewHandler(sessionService, userService),
		Util:     utilapi.NewHandler(keyGen),
		Health:   healthapi.NewHandler(),
	}

	middlewares := api.Middlewares{
		RequireActor:  middleware.RequireActor(cookieOperator, sessionService, localAdminEnabled),
		RequireAPIKey: middleware.RequireAPIKey(sessionAPIKey),
	}

	return handlers, middlewares
}
