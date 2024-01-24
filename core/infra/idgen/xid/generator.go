// Package xid provides a generator for IDs.
// Generated IDs are based on XIDs.
package xid

import (
	"github.com/energimind/identity-service/core/domain"
	"github.com/rs/xid"
)

// Generator implements the domain.IDGenerator interface, based on XID.
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
	return xid.New().String()
}
