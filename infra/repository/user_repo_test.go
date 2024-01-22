package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/identity-service/domain/auth"
	"github.com/energimind/identity-service/infra/repository"
	"github.com/energimind/identity-service/test/utils"
)

func TestUserRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewUserRepository(db)
	appID := auth.ID("1")

	utils.RunCRUDTests(t, utils.CRUDSetup[auth.User, auth.ID]{
		GetAll: func(ctx context.Context) ([]auth.User, error) {
			return repo.GetUsers(ctx, appID)
		},
		GetByID: func(ctx context.Context, id auth.ID) (auth.User, error) {
			return repo.GetUser(ctx, appID, id)
		},
		Create: func(ctx context.Context, user auth.User) error {
			return repo.CreateUser(ctx, user)
		},
		Update: func(ctx context.Context, user auth.User) error {
			return repo.UpdateUser(ctx, user)
		},
		Delete: func(ctx context.Context, id auth.ID) error {
			return repo.DeleteUser(ctx, appID, id)
		},
		NewEntity: func(key int) auth.User {
			return auth.User{
				ID:            auth.ID(strconv.Itoa(key)),
				ApplicationID: appID,
				Username:      "user1",
				Description:   "description",
				Enabled:       true,
				Role:          auth.SystemRoleAdmin,
				Accounts:      []auth.Account{{}},
				APIKeys:       []auth.APIKey{{}},
			}
		},
		ModifyEntity: func(user auth.User) auth.User {
			user.Username = "user2"

			return user
		},
		UnboundEntity: func() auth.User {
			return auth.User{
				ID:   "",
				Role: auth.SystemRoleAdmin,
			}
		},
		ExtractKey: func(user auth.User) auth.ID {
			return user.ID
		},
		MissingKey: func() auth.ID {
			return "missing"
		},
	})
}
