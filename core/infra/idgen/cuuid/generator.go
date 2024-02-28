// Package cuuid provides a generator for IDs.
// Generated IDs are condensed UUIDs.
package cuuid

import (
	"encoding/base32"

	"github.com/energimind/identity-server/core/domain"
	"github.com/google/uuid"
)

// Generator implements the domain.IDGenerator interface, based on XID.
type Generator struct {
	enc *base32.Encoding
}

// Ensure Generator implements the domain.IDGenerator interface.
var _ domain.IDGenerator = (*Generator)(nil)

// NewGenerator returns a new Generator instance.
func NewGenerator() *Generator {
	return &Generator{
		enc: base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding),
	}
}

// GenerateID generates a new ID.
//
// This method implements the idgen.Generator interface.
func (g *Generator) GenerateID() string {
	u := uuid.New()

	return g.enc.EncodeToString(u[:])
}
