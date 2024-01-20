package shortid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator_GenerateID(t *testing.T) {
	g := NewGenerator()

	id1 := g.GenerateID()
	id2 := g.GenerateID()

	require.NotEqual(t, id1, id2)
}
