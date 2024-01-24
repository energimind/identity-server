package xid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	t.Parallel()

	g := NewGenerator()

	id1 := g.GenerateID()
	id2 := g.GenerateID()
	id3 := g.GenerateID()

	require.Less(t, id1, id2)
	require.Less(t, id2, id3)
}
