package admin

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetUsers(t *testing.T) {
	t.Parallel()

	appID := admin.ID("a1")

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
			actor:      admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
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
			res, err := svc.GetUsers(context.Background(), test.actor, appID)

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

	appID := admin.ID("a1")
	userID := admin.ID("u1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:      admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantResult: true,
		},
		"user-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: "wrongAppID", UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"user-wrongUserID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: "wrongUserID"},
			wantError: domain.AccessDeniedError{},
		},
		"user-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantError: domain.StoreError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
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
			res, err := svc.GetUser(context.Background(), test.actor, appID, userID)

			if test.wantResult {
				require.NotEmpty(t, res)
			}

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	appID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
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
	svc := NewUserService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := admin.User{ApplicationID: appID}

			res, err := svc.CreateUser(context.Background(), test.actor, user)

			if test.wantResult {
				require.NotEmpty(t, res)
				require.NotEmpty(t, res.ID)
			}

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Parallel()

	userID := admin.ID("u1")
	appID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:      admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantResult: true,
		},
		"user-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: "wrongAppID", UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"user-wrongUserID": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: "wrongUserID"},
			wantError: domain.AccessDeniedError{},
		},
		"user-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantError: domain.StoreError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
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
			user := admin.User{ID: userID, ApplicationID: appID}
			res, err := svc.UpdateUser(context.Background(), test.actor, user)

			if test.wantResult {
				require.NotEmpty(t, res)
			}

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel()

	userID := admin.ID("u1")
	appID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
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
			err := svc.DeleteUser(context.Background(), test.actor, appID, userID)

			if test.wantResult {
				require.NoError(t, err)
			}

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockUserRepository struct {
	forcedError error
}

// ensure mockUserRepository implements admin.UserRepository.
var _ admin.UserRepository = (*mockUserRepository)(nil)

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{}
}

func (r *mockUserRepository) GetUsers(_ context.Context, appID admin.ID) ([]admin.User, error) {
	if appID == "" {
		return nil, errors.New("test-precondition: empty appID")
	}

	return []admin.User{r.mockUser()}, r.forcedError
}

func (r *mockUserRepository) GetUser(_ context.Context, appID, id admin.ID) (admin.User, error) {
	if appID == "" {
		return admin.User{}, errors.New("test-precondition: empty appID")
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

func (r *mockUserRepository) DeleteUser(_ context.Context, appID, id admin.ID) error {
	if appID == "" {
		return errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockUserRepository) GetUserByEmail(_ context.Context, appID admin.ID, email string) (admin.User, error) {
	if appID == "" {
		return admin.User{}, errors.New("test-precondition: empty appID")
	}

	if email == "" {
		return admin.User{}, errors.New("test-precondition: empty email")
	}

	return r.mockUser(), r.forcedError
}

func (r *mockUserRepository) mockUser() admin.User {
	return admin.User{
		ID:            "u1",
		ApplicationID: "a1",
		Username:      "mockUser",
	}
}
