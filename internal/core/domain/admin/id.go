package admin

// ID represents a unique identifier. It is a string type to allow for
// different types of identifiers.
type ID string

// IsEmpty returns true if the ID is empty.
func (i ID) IsEmpty() bool {
	return i == ""
}

// String returns the string representation of the ID.
func (i ID) String() string {
	return string(i)
}
