package cuuid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	t.Parallel()

	g := NewGenerator()

	prevID := ""

	for range 100 {
		id := g.GenerateID()

		require.NotEqual(t, prevID, id)
		require.Len(t, id, 26)
		require.NotContains(t, id, "-")

		prevID = id
	}
}
