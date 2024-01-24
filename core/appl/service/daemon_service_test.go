package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/stretchr/testify/require"
)

func TestDaemonService_GetDaemons(t *testing.T) {
	t.Parallel()

	appID := auth.ID("a1")

	tests := map[string]struct {
		actor      auth.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     auth.Actor{Role: auth.SystemRoleUser},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      auth.Actor{Role: auth.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     auth.Actor{Role: auth.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     auth.Actor{Role: "unknown"},
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

	appID := auth.ID("a1")
	userID := auth.ID("u1")

	tests := map[string]struct {
		actor      auth.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     auth.Actor{Role: auth.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      auth.Actor{Role: auth.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     auth.Actor{Role: auth.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     auth.Actor{Role: "unknown"},
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

	appID := auth.ID("a1")

	tests := map[string]struct {
		actor      auth.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     auth.Actor{Role: auth.SystemRoleUser, ApplicationID: appID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      auth.Actor{Role: auth.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     auth.Actor{Role: auth.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     auth.Actor{Role: "unknown"},
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
			user := auth.Daemon{ApplicationID: appID}

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

	userID := auth.ID("u1")
	appID := auth.ID("a1")

	tests := map[string]struct {
		actor      auth.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     auth.Actor{Role: auth.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      auth.Actor{Role: auth.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     auth.Actor{Role: auth.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     auth.Actor{Role: "unknown"},
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
			user := auth.Daemon{ID: userID, ApplicationID: appID}
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

	userID := auth.ID("u1")
	appID := auth.ID("a1")

	tests := map[string]struct {
		actor      auth.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     auth.Actor{Role: auth.SystemRoleUser, ApplicationID: appID, UserID: userID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:      auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantResult: true,
		},
		"manager-wrongAppID": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: "wrongAppID"},
			wantError: domain.AccessDeniedError{},
		},
		"manager-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.StoreError{},
		},
		"admin": {
			actor:      auth.Actor{Role: auth.SystemRoleAdmin},
			wantResult: true,
		},
		"admin-repoError": {
			actor:     auth.Actor{Role: auth.SystemRoleAdmin},
			wantError: domain.StoreError{},
		},
		"none": {
			actor:     auth.Actor{Role: auth.SystemRoleNone},
			wantError: domain.AccessDeniedError{},
		},
		"unknown": {
			actor:     auth.Actor{Role: "unknown"},
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

// ensure mockDaemonRepository implements auth.DaemonRepository.
var _ auth.DaemonRepository = (*mockDaemonRepository)(nil)

func newMockDaemonRepository() *mockDaemonRepository {
	return &mockDaemonRepository{}
}

func (r *mockDaemonRepository) GetDaemons(_ context.Context, appID auth.ID) ([]auth.Daemon, error) {
	if appID == "" {
		return nil, errors.New("test-precondition: empty appID")
	}

	return []auth.Daemon{r.mockDaemon()}, r.forcedError
}

func (r *mockDaemonRepository) GetDaemon(_ context.Context, appID, id auth.ID) (auth.Daemon, error) {
	if appID == "" {
		return auth.Daemon{}, errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return auth.Daemon{}, errors.New("test-precondition: empty id")
	}

	return r.mockDaemon(), r.forcedError
}

func (r *mockDaemonRepository) CreateDaemon(_ context.Context, user auth.Daemon) error {
	if (reflect.DeepEqual(user, auth.Daemon{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockDaemonRepository) UpdateDaemon(_ context.Context, user auth.Daemon) error {
	if (reflect.DeepEqual(user, auth.Daemon{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockDaemonRepository) DeleteDaemon(_ context.Context, appID, id auth.ID) error {
	if appID == "" {
		return errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockDaemonRepository) mockDaemon() auth.Daemon {
	return auth.Daemon{
		ID:            "u1",
		ApplicationID: "a1",
		Name:          "mockDaemon",
	}
}
