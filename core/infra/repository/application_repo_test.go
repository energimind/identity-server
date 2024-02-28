package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/infra/repository"
	"github.com/energimind/identity-server/core/test/utils"
)

func TestApplicationRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewApplicationRepository(db)

	utils.RunCRUDTests(t, utils.CRUDSetup[admin.Application, admin.ID]{
		GetAll: func(ctx context.Context) ([]admin.Application, error) {
			return repo.GetApplications(ctx)
		},
		GetByID: func(ctx context.Context, id admin.ID) (admin.Application, error) {
			return repo.GetApplication(ctx, id)
		},
		Create: func(ctx context.Context, app admin.Application) error {
			return repo.CreateApplication(ctx, app)
		},
		Update: func(ctx context.Context, app admin.Application) error {
			return repo.UpdateApplication(ctx, app)
		},
		Delete: func(ctx context.Context, id admin.ID) error {
			return repo.DeleteApplication(ctx, id)
		},
		NewEntity: func(key int) admin.Application {
			return admin.Application{
				ID:          admin.ID(strconv.Itoa(key)),
				Code:        "app1",
				Name:        "Application 1",
				Description: "Application 1",
				Enabled:     true,
			}
		},
		ModifyEntity: func(app admin.Application) admin.Application {
			app.Name = "Application 2"

			return app
		},
		UnboundEntity: func() admin.Application {
			return admin.Application{
				ID: "",
			}
		},
		ExtractKey: func(app admin.Application) admin.ID {
			return app.ID
		},
		MissingKey: func() admin.ID {
			return "missingKey"
		},
	})
}
