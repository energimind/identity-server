package server

import (
	"github.com/energimind/identity-service/core/api"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/energimind/identity-service/core/infra/rest/middleware"
)

func setupMiddlewares(cookieParser admin.CookieParser) api.Middlewares {
	return api.Middlewares{
		RequireActor: middleware.RequireActor(cookieParser),
	}
}
