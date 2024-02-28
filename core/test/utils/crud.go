package utils

import (
	"context"
	"testing"

	"github.com/energimind/identity-server/core/domain"
	"github.com/stretchr/testify/require"
)

// CRUDSetup is the setup for the CRUD tests.
type CRUDSetup[T, K any] struct {
	GetAll        func(ctx context.Context) ([]T, error)
	GetByID       func(ctx context.Context, id K) (T, error)
	Create        func(ctx context.Context, model T) error
	Update        func(ctx context.Context, model T) error
	Delete        func(ctx context.Context, id K) error
	NewEntity     func(key int) T
	ModifyEntity  func(model T) T
	UnboundEntity func() T
	ExtractKey    func(T) K
	MissingKey    func() K
}

// RunCRUDTests runs the CRUD tests for the given type.
func RunCRUDTests[T, K any](t *testing.T, setup CRUDSetup[T, K]) { //nolint:funlen
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyCounter := 0

	nextKey := func() int {
		keyCounter++

		return keyCounter
	}

	t.Run("missing-notFound", func(t *testing.T) {
		t.Run("getByID", func(t *testing.T) {
			_, err := setup.GetByID(ctx, setup.MissingKey())

			require.ErrorAs(t, err, &domain.NotFoundError{})
		})

		t.Run("update", func(t *testing.T) {
			require.ErrorAs(t, setup.Update(ctx, setup.UnboundEntity()), &domain.NotFoundError{})
		})

		t.Run("delete", func(t *testing.T) {
			require.ErrorAs(t, setup.Delete(ctx, setup.MissingKey()), &domain.NotFoundError{})
		})
	})

	t.Run("getAll-empty", func(t *testing.T) {
		all1, err := setup.GetAll(ctx)
		require.NoError(t, err)
		require.Empty(t, all1)
	})

	ent := setup.NewEntity(nextKey())
	key := setup.ExtractKey(ent)

	t.Run("create", func(t *testing.T) {
		require.NoError(t, setup.Create(ctx, ent))
	})

	t.Run("getAll-foundOne", func(t *testing.T) {
		all1, err := setup.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, []T{ent}, all1)
	})

	t.Run("getByID-found", func(t *testing.T) {
		e2, err := setup.GetByID(ctx, key)
		require.NoError(t, err)
		require.Equal(t, ent, e2)
	})

	entMod := setup.ModifyEntity(ent)

	require.Equal(t, key, setup.ExtractKey(entMod), "key should not change")
	require.NotEqual(t, ent, entMod, "entity should change")

	t.Run("update", func(t *testing.T) {
		require.NoError(t, setup.Update(ctx, entMod))
	})

	t.Run("getByID-updated", func(t *testing.T) {
		fetched, err := setup.GetByID(ctx, key)
		require.NoError(t, err)
		require.Equal(t, entMod, fetched)
	})

	t.Run("delete", func(t *testing.T) {
		require.NoError(t, setup.Delete(ctx, key))
	})

	t.Run("delete-again-notFound", func(t *testing.T) {
		require.ErrorAs(t, setup.Delete(ctx, key), &domain.NotFoundError{})
	})

	t.Run("getAll-empty", func(t *testing.T) {
		all2, err := setup.GetAll(ctx)
		require.NoError(t, err)
		require.Empty(t, all2)
	})
}
