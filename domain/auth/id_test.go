package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestID_IsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		id := ID("")

		require.True(t, id.IsEmpty())
	})

	t.Run("notEmpty", func(t *testing.T) {
		id := ID("1")

		require.False(t, id.IsEmpty())
	})
}
