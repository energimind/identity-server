package service

import (
	"context"
	"errors"
	"testing"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/stretchr/testify/require"
)

func TestApplicationService_GetApplications(t *testing.T) {
	t.Parallel()

	appID := auth.ID("1")

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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetApplications(context.Background(), test.actor)

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

func TestApplicationService_GetApplication(t *testing.T) {
	t.Parallel()

	appID := auth.ID("1")

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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, nil)

	id := appID

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetApplication(context.Background(), test.actor, id)

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

func TestApplicationService_CreateApplication(t *testing.T) {
	t.Parallel()

	appID := auth.ID("1")

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
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.AccessDeniedError{},
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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.CreateApplication(context.Background(), test.actor, auth.Application{})

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

func TestApplicationService_UpdateApplication(t *testing.T) {
	t.Parallel()

	appID := auth.ID("1")

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
			actor: auth.Actor{Role: auth.SystemRoleAdmin},
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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			app := auth.Application{ID: appID}
			res, err := svc.UpdateApplication(context.Background(), test.actor, app)

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

func TestApplicationService_DeleteApplication(t *testing.T) {
	t.Parallel()

	appID := auth.ID("1")

	tests := map[string]struct {
		actor     auth.Actor
		wantError error
	}{
		"user": {
			actor:     auth.Actor{Role: auth.SystemRoleUser},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:     auth.Actor{Role: auth.SystemRoleManager, ApplicationID: appID},
			wantError: domain.AccessDeniedError{},
		},
		"admin": {
			actor: auth.Actor{Role: auth.SystemRoleAdmin},
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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, nil)

	id := appID

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			err := svc.DeleteApplication(context.Background(), test.actor, id)

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockApplicationRepository struct {
	forcedError error
}

// ensure mockApplicationRepository implements auth.ApplicationRepository.
var _ auth.ApplicationRepository = (*mockApplicationRepository)(nil)

func newMockApplicationRepository() *mockApplicationRepository {
	return &mockApplicationRepository{}
}

func (r *mockApplicationRepository) GetApplications(_ context.Context) ([]auth.Application, error) {
	return []auth.Application{r.mockApplication()}, r.forcedError
}

func (r *mockApplicationRepository) GetApplication(_ context.Context, id auth.ID) (auth.Application, error) {
	if id == "" {
		return auth.Application{}, errors.New("test-precondition: empty id")
	}

	return r.mockApplication(), r.forcedError
}

func (r *mockApplicationRepository) CreateApplication(_ context.Context, app auth.Application) error {
	if (app == auth.Application{}) {
		return errors.New("test-precondition: empty application")
	}

	return r.forcedError
}

func (r *mockApplicationRepository) UpdateApplication(_ context.Context, app auth.Application) error {
	if (app == auth.Application{}) {
		return errors.New("test-precondition: empty application")
	}

	return r.forcedError
}

func (r *mockApplicationRepository) DeleteApplication(_ context.Context, id auth.ID) error {
	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockApplicationRepository) mockApplication() auth.Application {
	return auth.Application{
		ID:   "1",
		Name: "mockApplication",
	}
}
