package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/stretchr/testify/require"
)

func TestApplicationService_GetApplications(t *testing.T) {
	t.Parallel()

	appID := admin.ID("1")

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

	appID := admin.ID("1")

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

	appID := admin.ID("1")

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
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantError: domain.AccessDeniedError{},
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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.CreateApplication(context.Background(), test.actor, admin.Application{
				Code: "code",
				Name: "name",
			})

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

	appID := admin.ID("1")

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
			actor: admin.Actor{Role: admin.SystemRoleAdmin},
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

	repo := newMockApplicationRepository()
	svc := NewApplicationService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			app := admin.Application{
				ID:   appID,
				Code: "newCode",
				Name: "newName",
			}

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

	appID := admin.ID("1")

	tests := map[string]struct {
		actor     admin.Actor
		wantError error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, ApplicationID: appID},
			wantError: domain.AccessDeniedError{},
		},
		"admin": {
			actor: admin.Actor{Role: admin.SystemRoleAdmin},
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

// ensure mockApplicationRepository implements admin.ApplicationRepository.
var _ admin.ApplicationRepository = (*mockApplicationRepository)(nil)

func newMockApplicationRepository() *mockApplicationRepository {
	return &mockApplicationRepository{}
}

func (r *mockApplicationRepository) GetApplications(_ context.Context) ([]admin.Application, error) {
	return []admin.Application{r.mockApplication()}, r.forcedError
}

func (r *mockApplicationRepository) GetApplication(_ context.Context, id admin.ID) (admin.Application, error) {
	if id == "" {
		return admin.Application{}, errors.New("test-precondition: empty id")
	}

	return r.mockApplication(), r.forcedError
}

func (r *mockApplicationRepository) CreateApplication(_ context.Context, app admin.Application) error {
	if (app == admin.Application{}) {
		return errors.New("test-precondition: empty application")
	}

	return r.forcedError
}

func (r *mockApplicationRepository) UpdateApplication(_ context.Context, app admin.Application) error {
	if (app == admin.Application{}) {
		return errors.New("test-precondition: empty application")
	}

	return r.forcedError
}

func (r *mockApplicationRepository) DeleteApplication(_ context.Context, id admin.ID) error {
	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockApplicationRepository) mockApplication() admin.Application {
	return admin.Application{
		ID:   "1",
		Name: "mockApplication",
	}
}
