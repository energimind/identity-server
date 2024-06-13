package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetUsers(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantResult: true,
		},
		"manager-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: "wrongRealmID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      admin.Actor{Role: admin.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     admin.Actor{Role: admin.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     admin.Actor{Role: "unknown"},
			wantError: domain.AccessDeniedError{},
		},
	}

	repo := newMockUserRepository()
	svc := NewUserService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetUsers(context.Background(), test.actor, realmID)

			if test.wantResult {
				require.Len(t, res, 1)
			}

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("a1")
	userID := admin.ID("u1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:      admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: userID},
			wantResult: true,
		},
		"user-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: "wrongRealmID", UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"user-wrongUserID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: "wrongUserID"},
			wantError: domain.AccessDeniedError{},
		},
		"user-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: userID},
			wantError: domain.StoreError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantResult: true,
		},
		"manager-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: "wrongRealmID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      admin.Actor{Role: admin.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     admin.Actor{Role: admin.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     admin.Actor{Role: "unknown"},
			wantError: domain.AccessDeniedError{},
		},
	}

	repo := newMockUserRepository()
	svc := NewUserService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetUser(context.Background(), test.actor, realmID, userID)

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}

			if test.wantResult {
				require.NotEmpty(t, res)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		userExists bool
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantResult: true,
		},
		"manager-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: "wrongRealmID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantError: domain.StoreError{},
		},
		"manager-duplicateUser": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			userExists: true,
			wantError:  domain.ConflictError{},
		},
		"admin": {
			actor:      admin.Actor{Role: admin.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     admin.Actor{Role: admin.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     admin.Actor{Role: "unknown"},
			wantError: domain.AccessDeniedError{},
		},
	}

	repo := newMockUserRepository()
	svc := NewUserService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		repo.userExists = test.userExists

		t.Run(name, func(t *testing.T) {
			user := admin.User{
				RealmID:  realmID,
				BindID:   "bindID",
				Username: "testUser",
				Email:    "email@domain.com",
			}

			res, err := svc.CreateUser(context.Background(), test.actor, user)

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}

			if test.wantResult {
				require.NotEmpty(t, res)
				require.NotEmpty(t, res.ID)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Parallel()

	userID := admin.ID("u1")
	realmID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:      admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: userID},
			wantResult: true,
		},
		"user-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: "wrongRealmID", UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"user-wrongUserID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: "wrongUserID"},
			wantError: domain.AccessDeniedError{},
		},
		"user-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: userID},
			wantError: domain.StoreError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantResult: true,
		},
		"manager-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: "wrongRealmID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      admin.Actor{Role: admin.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     admin.Actor{Role: admin.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     admin.Actor{Role: "unknown"},
			wantError: domain.AccessDeniedError{},
		},
	}

	repo := newMockUserRepository()
	svc := NewUserService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := admin.User{
				ID:       userID,
				RealmID:  realmID,
				BindID:   "bindID",
				Username: "newUsername",
				Email:    "newMail@domain.com",
			}

			res, err := svc.UpdateUser(context.Background(), test.actor, user)

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}

			if test.wantResult {
				require.NotEmpty(t, res)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel()

	userID := admin.ID("u1")
	realmID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID, UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantResult: true,
		},
		"manager-wrongRealmID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: "wrongRealmID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      admin.Actor{Role: admin.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     admin.Actor{Role: admin.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     admin.Actor{Role: "unknown"},
			wantError: domain.AccessDeniedError{},
		},
	}

	repo := newMockUserRepository()
	svc := NewUserService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			err := svc.DeleteUser(context.Background(), test.actor, realmID, userID)

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}

			if test.wantResult {
				require.NoError(t, err)
			}
		})
	}
}

type mockUserRepository struct {
	userExists  bool
	forcedError error
}

// ensure mockUserRepository implements admin.UserRepository.
var _ admin.UserRepository = (*mockUserRepository)(nil)

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{}
}

func (r *mockUserRepository) GetUsers(_ context.Context, realmID admin.ID) ([]admin.User, error) {
	if realmID == "" {
		return nil, errors.New("test-precondition: empty realmID")
	}

	return []admin.User{r.mockUser()}, r.forcedError
}

func (r *mockUserRepository) GetUser(_ context.Context, realmID, id admin.ID) (admin.User, error) {
	if realmID == "" {
		return admin.User{}, errors.New("test-precondition: empty realmID")
	}

	if id == "" {
		return admin.User{}, errors.New("test-precondition: empty id")
	}

	return r.mockUser(), r.forcedError
}

func (r *mockUserRepository) CreateUser(_ context.Context, user admin.User) error {
	if (reflect.DeepEqual(user, admin.User{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockUserRepository) UpdateUser(_ context.Context, user admin.User) error {
	if (reflect.DeepEqual(user, admin.User{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockUserRepository) DeleteUser(_ context.Context, realmID, id admin.ID) error {
	if realmID == "" {
		return errors.New("test-precondition: empty realmID")
	}

	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockUserRepository) GetUserByBindID(_ context.Context, realmID admin.ID, bindID string) (admin.User, error) {
	if realmID == "" {
		return admin.User{}, errors.New("test-precondition: empty realmID")
	}

	if bindID == "" {
		return admin.User{}, errors.New("test-precondition: empty bindID")
	}

	if !r.userExists {
		return admin.User{}, domain.NewNotFoundError("user not found")
	}

	return r.mockUser(), r.forcedError
}

func (r *mockUserRepository) GetAPIKey(_ context.Context, realmID admin.ID, key string) (admin.APIKey, error) {
	if realmID == "" {
		return admin.APIKey{}, errors.New("test-precondition: empty realmID")
	}

	if key == "" {
		return admin.APIKey{}, errors.New("test-precondition: empty key")
	}

	return admin.APIKey{}, r.forcedError
}

func (r *mockUserRepository) mockUser() admin.User {
	return admin.User{
		ID:       "u1",
		RealmID:  "a1",
		Username: "mockUser",
	}
}
