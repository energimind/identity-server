package server

import (
	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/rest/middleware"
)

func setupMiddlewares(cookieVerifier auth.CookieVerifier) api.Middlewares {
	return api.Middlewares{
		RequireActor: middleware.RequireActor(cookieVerifier),
	}
}
