package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/identity-service/domain/auth"
	"github.com/energimind/identity-service/infra/repository"
	"github.com/energimind/identity-service/test/utils"
)

func TestProviderRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewProviderRepository(db)

	utils.RunCRUDTests(t, utils.CRUDSetup[auth.Provider, auth.ID]{
		GetAll: func(ctx context.Context) ([]auth.Provider, error) {
			return repo.GetProviders(ctx)
		},
		GetByID: func(ctx context.Context, id auth.ID) (auth.Provider, error) {
			return repo.GetProvider(ctx, id)
		},
		Create: func(ctx context.Context, provider auth.Provider) error {
			return repo.CreateProvider(ctx, provider)
		},
		Update: func(ctx context.Context, provider auth.Provider) error {
			return repo.UpdateProvider(ctx, provider)
		},
		Delete: func(ctx context.Context, id auth.ID) error {
			return repo.DeleteProvider(ctx, id)
		},
		NewEntity: func(key int) auth.Provider {
			return auth.Provider{
				ID:            auth.ID(strconv.Itoa(key)),
				ApplicationID: "app1",
				Type:          auth.ProviderTypeGoogle,
				Code:          "google",
				Name:          "Google",
				Description:   "Google",
				Enabled:       true,
				ClientID:      "client-id",
				ClientSecret:  "client-secret",
				RedirectURL:   "https://example.com",
			}
		},
		ModifyEntity: func(provider auth.Provider) auth.Provider {
			provider.Name = "Google 2"

			return provider
		},
		UnboundEntity: func() auth.Provider {
			return auth.Provider{
				ID:   "",
				Type: auth.ProviderTypeGoogle,
			}
		},
		ExtractKey: func(provider auth.Provider) auth.ID {
			return provider.ID
		},
		MissingKey: func() auth.ID {
			return "missing"
		},
	})
}
