package service

import (
	"context"
	"errors"
	"testing"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/stretchr/testify/require"
)

func TestRealmService_GetRealms(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("1")

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

	repo := newMockRealmRepository()
	svc := NewRealmService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetRealms(context.Background(), test.actor)

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

func TestRealmService_GetRealm(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("1")

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

	repo := newMockRealmRepository()
	svc := NewRealmService(repo, nil)

	id := realmID

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.GetRealm(context.Background(), test.actor, id)

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

func TestRealmService_CreateRealm(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("1")

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

	repo := newMockRealmRepository()
	svc := NewRealmService(repo, newMockIDGenerator())

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			res, err := svc.CreateRealm(context.Background(), test.actor, admin.Realm{
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

func TestRealmService_UpdateRealm(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("1")

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

	repo := newMockRealmRepository()
	svc := NewRealmService(repo, nil)

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			realm := admin.Realm{
				ID:   realmID,
				Code: "newCode",
				Name: "newName",
			}

			res, err := svc.UpdateRealm(context.Background(), test.actor, realm)

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

func TestRealmService_DeleteRealm(t *testing.T) {
	t.Parallel()

	realmID := admin.ID("1")

	tests := map[string]struct {
		actor     admin.Actor
		wantError error
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

	repo := newMockRealmRepository()
	svc := NewRealmService(repo, nil)

	id := realmID

	for name, test := range tests {
		if errors.Is(test.wantError, domain.StoreError{}) {
			repo.forcedError = errors.New("forcedError")
		} else {
			repo.forcedError = nil
		}

		t.Run(name, func(t *testing.T) {
			err := svc.DeleteRealm(context.Background(), test.actor, id)

			if test.wantError != nil {
				require.ErrorAs(t, err, &test.wantError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockRealmRepository struct {
	forcedError error
}

// ensure mockRealmRepository implements admin.RealmRepository.
var _ admin.RealmRepository = (*mockRealmRepository)(nil)

func newMockRealmRepository() *mockRealmRepository {
	return &mockRealmRepository{}
}

func (r *mockRealmRepository) GetRealms(_ context.Context) ([]admin.Realm, error) {
	return []admin.Realm{r.mockRealm()}, r.forcedError
}

func (r *mockRealmRepository) GetRealm(_ context.Context, id admin.ID) (admin.Realm, error) {
	if id == "" {
		return admin.Realm{}, errors.New("test-precondition: empty id")
	}

	return r.mockRealm(), r.forcedError
}

func (r *mockRealmRepository) CreateRealm(_ context.Context, realm admin.Realm) error {
	if (realm == admin.Realm{}) {
		return errors.New("test-precondition: empty realm")
	}

	return r.forcedError
}

func (r *mockRealmRepository) UpdateRealm(_ context.Context, realm admin.Realm) error {
	if (realm == admin.Realm{}) {
		return errors.New("test-precondition: empty realm")
	}

	return r.forcedError
}

func (r *mockRealmRepository) DeleteRealm(_ context.Context, id admin.ID) error {
	if id == "" {
		return errors.New("test-precondition: empty id")
	}

	return r.forcedError
}

func (r *mockRealmRepository) mockRealm() admin.Realm {
	return admin.Realm{
		ID:   "1",
		Name: "mockRealm",
	}
}
