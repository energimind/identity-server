package domain

// IDGenerator is an interface for generating IDs.
type IDGenerator interface {
	GenerateID() string
}
