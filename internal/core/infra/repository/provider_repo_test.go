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

func TestProviderRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewProviderRepository(db)

	crud.RunTests(t, crud.Setup[admin.Provider, admin.ID]{
		RepoOps: crud.RepoOps[admin.Provider, admin.ID]{
			GetAll: func(ctx context.Context) ([]admin.Provider, error) {
				return repo.GetProviders(ctx)
			},
			GetByID: func(ctx context.Context, id admin.ID) (admin.Provider, error) {
				return repo.GetProvider(ctx, id)
			},
			Create: func(ctx context.Context, provider admin.Provider) error {
				return repo.CreateProvider(ctx, provider)
			},
			Update: func(ctx context.Context, provider admin.Provider) error {
				return repo.UpdateProvider(ctx, provider)
			},
			Delete: func(ctx context.Context, id admin.ID) error {
				return repo.DeleteProvider(ctx, id)
			},
		},
		EntityOps: crud.EntityOps[admin.Provider, admin.ID]{
			NewEntity: func(key int) admin.Provider {
				return admin.Provider{
					ID:           admin.ID(strconv.Itoa(key)),
					Type:         admin.ProviderTypeGoogle,
					Code:         "google",
					Name:         "Google",
					Description:  "Google",
					Enabled:      true,
					ClientID:     "client-id",
					ClientSecret: "client-secret",
					RedirectURL:  "https://example.com",
				}
			},
			ModifyEntity: func(provider admin.Provider) admin.Provider {
				provider.Name = "Google 2"

				return provider
			},
			UnboundEntity: func() admin.Provider {
				return admin.Provider{
					ID:   "",
					Type: admin.ProviderTypeGoogle,
				}
			},
			ExtractKey: func(provider admin.Provider) admin.ID {
				return provider.ID
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
