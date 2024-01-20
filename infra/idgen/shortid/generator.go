// Package shortid provides a generator for IDs.
// Generated IDs are based on short IDs.
package shortid

import (
	"math/rand"
	"time"

	"github.com/energimind/identity-service/domain"
	"github.com/teris-io/shortid"
)

// Generator implements the domain.IDGenerator interface, based on shortid.
type Generator struct {
	gen *shortid.Shortid
}

// Ensure Generator implements the domain.IDGenerator interface.
var _ domain.IDGenerator = (*Generator)(nil)

// NewGenerator returns a new Generator instance.
func NewGenerator() *Generator {
	//nolint:gosec // this is good enough for our use case
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Generator{
		gen: shortid.MustNew(1, shortid.DefaultABC, r.Uint64()),
	}
}

// GenerateID generates a new ID.
//
// This method implements the idgen.Generator interface.
func (g *Generator) GenerateID() string {
	return g.gen.MustGenerate()
}
