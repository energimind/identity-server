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

func TestProviderService_GetProviders(t *testing.T) {
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
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetProviders(context.Background(), test.actor)

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

	realmID := admin.ID("a1")
	userID := admin.ID("u1")

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
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetProvider(context.Background(), test.actor, userID)

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

	realmID := admin.ID("a1")

	tests := map[string]struct {
		actor      admin.Actor
		wantResult bool
		wantError  error
	}{
		"user": {
			actor:     admin.Actor{Role: admin.SystemRoleUser, RealmID: realmID},
			wantError: domain.AccessDeniedError{},
		},
		"manager": {
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := admin.Provider{
				Code: "code",
				Name: "name",
			}

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
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			user := admin.Provider{
				ID:   userID,
				Code: "newCode",
				Name: "newName",
			}

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
			actor:     admin.Actor{Role: admin.SystemRoleManager, RealmID: realmID},
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

	repo := newMockProviderRepository()
	svc := NewProviderService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			err := svc.DeleteProvider(context.Background(), test.actor, userID)

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

// ensure mockProviderRepository implements admin.ProviderRepository.
var _ admin.ProviderRepository = (*mockProviderRepository)(nil)

func newMockProviderRepository() *mockProviderRepository {
	return &mockProviderRepository{}
}

func (r *mockProviderRepository) GetProviders(_ context.Context) ([]admin.Provider, error) {
	return []admin.Provider{r.mockProvider()}, r.forcedError
}

func (r *mockProviderRepository) GetProvider(_ context.Context, id admin.ID) (admin.Provider, error) {
	if id == "" {
		return admin.Provider{}, errors.New("test-precondition: empty id")
	}

	return r.mockProvider(), r.forcedError
}

func (r *mockProviderRepository) CreateProvider(_ context.Context, user admin.Provider) error {
	if (reflect.DeepEqual(user, admin.Provider{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockProviderRepository) UpdateProvider(_ context.Context, user admin.Provider) error {
	if (reflect.DeepEqual(user, admin.Provider{})) {
		return errors.New("test-precondition: empty user")
	}

	return r.forcedError
}

func (r *mockProviderRepository) DeleteProvider(_ context.Context, id admin.ID) error {
	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockProviderRepository) mockProvider() admin.Provider {
	return admin.Provider{
		ID:   "u1",
		Name: "mockProvider",
	}
}
