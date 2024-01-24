package repository

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_drainCursor(t *testing.T) {
	t.Parallel()

	mapper := func(s string) string {
		return strings.ToUpper(s)
	}

	t.Run("success", func(t *testing.T) {
		cur := &mockCursor{
			all: func(v any) error {
				*v.(*[]string) = []string{"a", "b"}
				return nil
			},
			close: func() error {
				return nil
			},
		}

		s, err := drainCursor[string](context.Background(), cur, mapper)

		require.NoError(t, err)
		require.Equal(t, []string{"A", "B"}, s)
	})

	t.Run("allError", func(t *testing.T) {
		cur := &mockCursor{
			all: func(v any) error {
				return errors.New("forcedError")
			},
		}

		res, err := drainCursor[string](context.Background(), cur, mapper)

		require.ErrorContains(t, err, "forcedError")
		require.Nil(t, res)
	})

	t.Run("closeError", func(t *testing.T) {
		cur := &mockCursor{
			all: func(v any) error {
				return nil
			},
			close: func() error {
				return errors.New("forcedError")
			},
		}

		res, err := drainCursor[string](context.Background(), cur, mapper)

		require.ErrorContains(t, err, "forcedError")
		require.Nil(t, res)
	})
}

type mockCursor struct {
	all   func(v any) error
	close func() error
}

// ensure mockCursor implements cursor interface
var _ cursor = (*mockCursor)(nil)

func (m *mockCursor) All(_ context.Context, v any) error {
	return m.all(v)
}

func (m *mockCursor) Close(_ context.Context) error {
	return m.close()
}
