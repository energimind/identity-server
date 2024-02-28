package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/infra/repository"
	"github.com/energimind/identity-server/core/test/utils"
)

func TestProviderRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewProviderRepository(db)
	appID := admin.ID("1")

	utils.RunCRUDTests(t, utils.CRUDSetup[admin.Provider, admin.ID]{
		GetAll: func(ctx context.Context) ([]admin.Provider, error) {
			return repo.GetProviders(ctx, appID)
		},
		GetByID: func(ctx context.Context, id admin.ID) (admin.Provider, error) {
			return repo.GetProvider(ctx, appID, id)
		},
		Create: func(ctx context.Context, provider admin.Provider) error {
			return repo.CreateProvider(ctx, provider)
		},
		Update: func(ctx context.Context, provider admin.Provider) error {
			return repo.UpdateProvider(ctx, provider)
		},
		Delete: func(ctx context.Context, id admin.ID) error {
			return repo.DeleteProvider(ctx, appID, id)
		},
		NewEntity: func(key int) admin.Provider {
			return admin.Provider{
				ID:            admin.ID(strconv.Itoa(key)),
				ApplicationID: appID,
				Type:          admin.ProviderTypeGoogle,
				Code:          "google",
				Name:          "Google",
				Description:   "Google",
				Enabled:       true,
				ClientID:      "client-id",
				ClientSecret:  "client-secret",
				RedirectURL:   "https://example.com",
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
	})
}
