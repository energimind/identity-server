package repository_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/energimind/go-kit/testutil/crud"
	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/infra/repository"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CRUD(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewUserRepository(db)
	realmID := admin.ID("1")

	crud.RunTests(t, crud.Setup[admin.User, admin.ID]{
		RepoOps: crud.RepoOps[admin.User, admin.ID]{
			GetAll: func(ctx context.Context) ([]admin.User, error) {
				return repo.GetUsers(ctx, realmID)
			},
			GetByID: func(ctx context.Context, id admin.ID) (admin.User, error) {
				return repo.GetUser(ctx, realmID, id)
			},
			Create: func(ctx context.Context, user admin.User) error {
				return repo.CreateUser(ctx, user)
			},
			Update: func(ctx context.Context, user admin.User) error {
				return repo.UpdateUser(ctx, user)
			},
			Delete: func(ctx context.Context, id admin.ID) error {
				return repo.DeleteUser(ctx, realmID, id)
			},
		},
		EntityOps: crud.EntityOps[admin.User, admin.ID]{
			NewEntity: func(key int) admin.User {
				return admin.User{
					ID:          admin.ID(strconv.Itoa(key)),
					RealmID:     realmID,
					Username:    "user1",
					Description: "description",
					Enabled:     true,
					Role:        admin.SystemRoleAdmin,
					APIKeys:     []admin.APIKey{{}},
				}
			},
			ModifyEntity: func(user admin.User) admin.User {
				user.Username = "user2"

				return user
			},
			UnboundEntity: func() admin.User {
				return admin.User{
					ID:   "",
					Role: admin.SystemRoleAdmin,
				}
			},
			ExtractKey: func(user admin.User) admin.ID {
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

func TestUserRepository_GetUserByEmail(t *testing.T) {
	t.Parallel()

	db, closer := mongoEnv.NewInstance()
	defer closer()

	repo := repository.NewUserRepository(db)
	realmID := admin.ID("1")

	ctx := context.Background()
	user := admin.User{
		ID:      "1",
		RealmID: realmID,
		Email:   "user@somedomain.com",
		APIKeys: []admin.APIKey{},
	}

	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	got, err := repo.GetUserByEmail(ctx, realmID, user.Email)
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	require.Equal(t, user, got)
}
