package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/go-kit/testutil/crud"
	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/infra/repository"
)

func TestDaemonRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewDaemonRepository(db)
	appID := admin.ID("1")

	crud.RunTests(t, crud.Setup[admin.Daemon, admin.ID]{
		RepoOps: crud.RepoOps[admin.Daemon, admin.ID]{
			GetAll: func(ctx context.Context) ([]admin.Daemon, error) {
				return repo.GetDaemons(ctx, appID)
			},
			GetByID: func(ctx context.Context, id admin.ID) (admin.Daemon, error) {
				return repo.GetDaemon(ctx, appID, id)
			},
			Create: func(ctx context.Context, user admin.Daemon) error {
				return repo.CreateDaemon(ctx, user)
			},
			Update: func(ctx context.Context, user admin.Daemon) error {
				return repo.UpdateDaemon(ctx, user)
			},
			Delete: func(ctx context.Context, id admin.ID) error {
				return repo.DeleteDaemon(ctx, appID, id)
			},
		},
		EntityOps: crud.EntityOps[admin.Daemon, admin.ID]{
			NewEntity: func(key int) admin.Daemon {
				return admin.Daemon{
					ID:            admin.ID(strconv.Itoa(key)),
					ApplicationID: appID,
					Code:          "daemon",
					Name:          "Daemon",
					Description:   "Daemon description",
					Enabled:       true,
					APIKeys:       []admin.APIKey{{}},
				}
			},
			ModifyEntity: func(user admin.Daemon) admin.Daemon {
				user.Name = "Daemon 2"

				return user
			},
			UnboundEntity: func() admin.Daemon {
				return admin.Daemon{ID: ""}
			},
			ExtractKey: func(user admin.Daemon) admin.ID {
				return user.ID
			},
			MissingKey: func() admin.ID {
				return "missing"
			},
		},
		NotFoundErr: func() any {
			return domain.NotFoundError{}
		},
	})
}
