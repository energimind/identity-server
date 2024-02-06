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

func TestDaemonService_GetDaemons(t *testing.T) {
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

	repo := newMockDaemonRepository()
	svc := NewDaemonService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetDaemons(context.Background(), test.actor, appID)

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

func TestDaemonService_GetDaemon(t *testing.T) {
	t.Parallel()

	appID := admin.ID("a1")
	userID := admin.ID("u1")

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

	repo := newMockDaemonRepository()
	svc := NewDaemonService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetDaemon(context.Background(), test.actor, appID, userID)

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

func TestDaemonService_CreateDaemon(t *testing.T) {
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

	repo := newMockDaemonRepository()
	svc := NewDaemonService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := admin.Daemon{ApplicationID: appID}

			res, err := svc.CreateDaemon(context.Background(), test.actor, user)

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

func TestDaemonService_UpdateDaemon(t *testing.T) {
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

	repo := newMockDaemonRepository()
	svc := NewDaemonService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := admin.Daemon{ID: userID, ApplicationID: appID}
			res, err := svc.UpdateDaemon(context.Background(), test.actor, user)

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

func TestDaemonService_DeleteDaemon(t *testing.T) {
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

	repo := newMockDaemonRepository()
	svc := NewDaemonService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			err := svc.DeleteDaemon(context.Background(), test.actor, appID, userID)

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

type mockDaemonRepository struct {
	forcedError error
}

// ensure mockDaemonRepository implements admin.DaemonRepository.
var _ admin.DaemonRepository = (*mockDaemonRepository)(nil)

func newMockDaemonRepository() *mockDaemonRepository {
	return &mockDaemonRepository{}
}

func (r *mockDaemonRepository) GetDaemons(_ context.Context, appID admin.ID) ([]admin.Daemon, error) {
	if appID == "" {
		return nil, errors.New("test-precondition: empty appID")
	}

	return []admin.Daemon{r.mockDaemon()}, r.forcedError
}

func (r *mockDaemonRepository) GetDaemon(_ context.Context, appID, id admin.ID) (admin.Daemon, error) {
	if appID == "" {
		return admin.Daemon{}, errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return admin.Daemon{}, errors.New("test-precondition: empty id")
	}

	return r.mockDaemon(), r.forcedError
}

func (r *mockDaemonRepository) CreateDaemon(_ context.Context, user admin.Daemon) error {
	if (reflect.DeepEqual(user, admin.Daemon{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockDaemonRepository) UpdateDaemon(_ context.Context, user admin.Daemon) error {
	if (reflect.DeepEqual(user, admin.Daemon{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockDaemonRepository) DeleteDaemon(_ context.Context, appID, id admin.ID) error {
	if appID == "" {
		return errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockDaemonRepository) mockDaemon() admin.Daemon {
	return admin.Daemon{
		ID:            "u1",
		ApplicationID: "a1",
		Name:          "mockDaemon",
	}
}
