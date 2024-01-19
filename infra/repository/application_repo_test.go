package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/identity-service/domain/auth"
	"github.com/energimind/identity-service/infra/repository"
	"github.com/energimind/identity-service/test/utils"
)

func TestApplicationRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewApplicationRepository(db)
	allowAll := auth.Principal{}

	utils.RunCRUDTests(t, utils.CRUDSetup[auth.Application, auth.ID]{
		GetAll: func(ctx context.Context) ([]auth.Application, error) {
			return repo.GetApplications(ctx, allowAll)
		},
		GetByID: func(ctx context.Context, id auth.ID) (auth.Application, error) {
			return repo.GetApplication(ctx, allowAll, id)
		},
		Create: func(ctx context.Context, app auth.Application) error {
			return repo.CreateApplication(ctx, allowAll, app)
		},
		Update: func(ctx context.Context, app auth.Application) error {
			return repo.UpdateApplication(ctx, allowAll, app)
		},
		Delete: func(ctx context.Context, id auth.ID) error {
			return repo.DeleteApplication(ctx, allowAll, id)
		},
		NewEntity: func(key int) auth.Application {
			return auth.Application{
				ID:          auth.ID(strconv.Itoa(key)),
				Code:        "app1",
				Name:        "Application 1",
				Description: "Application 1",
				Enabled:     true,
			}
		},
		ModifyEntity: func(app auth.Application) auth.Application {
			app.Name = "Application 2"

			return app
		},
		UnboundEntity: func() auth.Application {
			return auth.Application{
				ID: "",
			}
		},
		ExtractKey: func(app auth.Application) auth.ID {
			return app.ID
		},
		MissingKey: func() auth.ID {
			return "missingKey"
		},
	})
}
