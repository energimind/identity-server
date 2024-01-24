package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/energimind/identity-service/core/infra/repository"
	"github.com/energimind/identity-service/test/utils"
)

func TestDaemonRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewDaemonRepository(db)
	appID := auth.ID("1")

	utils.RunCRUDTests(t, utils.CRUDSetup[auth.Daemon, auth.ID]{
		GetAll: func(ctx context.Context) ([]auth.Daemon, error) {
			return repo.GetDaemons(ctx, appID)
		},
		GetByID: func(ctx context.Context, id auth.ID) (auth.Daemon, error) {
			return repo.GetDaemon(ctx, appID, id)
		},
		Create: func(ctx context.Context, user auth.Daemon) error {
			return repo.CreateDaemon(ctx, user)
		},
		Update: func(ctx context.Context, user auth.Daemon) error {
			return repo.UpdateDaemon(ctx, user)
		},
		Delete: func(ctx context.Context, id auth.ID) error {
			return repo.DeleteDaemon(ctx, appID, id)
		},
		NewEntity: func(key int) auth.Daemon {
			return auth.Daemon{
				ID:            auth.ID(strconv.Itoa(key)),
				ApplicationID: appID,
				Code:          "daemon",
				Name:          "Daemon",
				Description:   "Daemon description",
				Enabled:       true,
				APIKeys:       []auth.APIKey{{}},
			}
		},
		ModifyEntity: func(user auth.Daemon) auth.Daemon {
			user.Name = "Daemon 2"

			return user
		},
		UnboundEntity: func() auth.Daemon {
			return auth.Daemon{ID: ""}
		},
		ExtractKey: func(user auth.Daemon) auth.ID {
			return user.ID
		},
		MissingKey: func() auth.ID {
			return "missing"
		},
	})
}
