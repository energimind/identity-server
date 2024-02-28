// Package uuid provides a generator for IDs based on UUIDs.
package uuid

import (
	"github.com/energimind/identity-server/core/domain"
	"github.com/google/uuid"
)

// Generator implements the domain.IDGenerator interface, based on UUIDs.
type Generator struct{}

// Ensure Generator implements the domain.IDGenerator interface.
var _ domain.IDGenerator = (*Generator)(nil)

// NewGenerator returns a new Generator instance.
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateID generates a new ID.
//
// This method implements the idgen.Generator interface.
func (g *Generator) GenerateID() string {
	return uuid.New().String()
}
