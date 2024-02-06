package server

import (
	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/energimind/identity-service/core/infra/rest/middleware"
)

func setupMiddlewares(cookieVerifier admin.CookieVerifier) api.Middlewares {
	return api.Middlewares{
		RequireActor: middleware.RequireActor(cookieVerifier),
	}
}
