package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
	"github.com/stretchr/testify/require"
)

func TestProviderService_GetProviders(t *testing.T) {
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetProviders(context.Background(), test.actor, appID)

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

func TestProviderService_GetProvider(t *testing.T) {
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetProvider(context.Background(), test.actor, appID, userID)

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

func TestProviderService_CreateProvider(t *testing.T) {
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := auth.Provider{ApplicationID: appID}

			res, err := svc.CreateProvider(context.Background(), test.actor, user)

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

func TestProviderService_UpdateProvider(t *testing.T) {
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := auth.Provider{ID: userID, ApplicationID: appID}
			res, err := svc.UpdateProvider(context.Background(), test.actor, user)

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

func TestProviderService_DeleteProvider(t *testing.T) {
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			err := svc.DeleteProvider(context.Background(), test.actor, appID, userID)

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

type mockProviderRepository struct {
	forcedError error
}

// ensure mockProviderRepository implements auth.ProviderRepository.
var _ auth.ProviderRepository = (*mockProviderRepository)(nil)

func newMockProviderRepository() *mockProviderRepository {
	return &mockProviderRepository{}
}

func (r *mockProviderRepository) GetProviders(_ context.Context, appID auth.ID) ([]auth.Provider, error) {
	if appID == "" {
		return nil, errors.New("test-precondition: empty appID")
	}

	return []auth.Provider{r.mockProvider()}, r.forcedError
}

func (r *mockProviderRepository) GetProvider(_ context.Context, appID, id auth.ID) (auth.Provider, error) {
	if appID == "" {
		return auth.Provider{}, errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return auth.Provider{}, errors.New("test-precondition: empty id")
	}

	return r.mockProvider(), r.forcedError
}

func (r *mockProviderRepository) CreateProvider(_ context.Context, user auth.Provider) error {
	if (reflect.DeepEqual(user, auth.Provider{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockProviderRepository) UpdateProvider(_ context.Context, user auth.Provider) error {
	if (reflect.DeepEqual(user, auth.Provider{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockProviderRepository) DeleteProvider(_ context.Context, appID, id auth.ID) error {
	if appID == "" {
		return errors.New("test-precondition: empty appID")
	}

	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockProviderRepository) mockProvider() auth.Provider {
	return auth.Provider{
		ID:            "u1",
		ApplicationID: "a1",
		Name:          "mockProvider",
	}
}
