package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/go-kit/testutil/crud"
	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/infra/repository"
)

func TestRealmRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewRealmRepository(db)

	crud.RunTests(t, crud.Setup[admin.Realm, admin.ID]{
		RepoOps: crud.RepoOps[admin.Realm, admin.ID]{
			GetAll: func(ctx context.Context) ([]admin.Realm, error) {
				return repo.GetRealms(ctx)
			},
			GetByID: func(ctx context.Context, id admin.ID) (admin.Realm, error) {
				return repo.GetRealm(ctx, id)
			},
			Create: func(ctx context.Context, realm admin.Realm) error {
				return repo.CreateRealm(ctx, realm)
			},
			Update: func(ctx context.Context, realm admin.Realm) error {
				return repo.UpdateRealm(ctx, realm)
			},
			Delete: func(ctx context.Context, id admin.ID) error {
				return repo.DeleteRealm(ctx, id)
			},
		},
		EntityOps: crud.EntityOps[admin.Realm, admin.ID]{
			NewEntity: func(key int) admin.Realm {
				return admin.Realm{
					ID:          admin.ID(strconv.Itoa(key)),
					Code:        "realm1",
					Name:        "Realm 1",
					Description: "Realm 1",
					Enabled:     true,
				}
			},
			ModifyEntity: func(realm admin.Realm) admin.Realm {
				realm.Name = "Realm 2"

				return realm
			},
			UnboundEntity: func() admin.Realm {
				return admin.Realm{
					ID: "",
				}
			},
			ExtractKey: func(realm admin.Realm) admin.ID {
				return realm.ID
			},
			MissingKey: func() admin.ID {
				return "missingKey"
			},
		},
		NotFoundErr: func() any {
			return domain.NotFoundError{}
		},
	})
}
