package redis

import "fmt"

// CacheError represents a cache error.
type CacheError struct {
	Message string
}

// NewCacheError creates a new CacheError.
func NewCacheError(format string, args ...any) CacheError {
	return CacheError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error implements the error interface.
func (e CacheError) Error() string {
	return e.Message
}

// String implements the fmt.Stringer interface.
func (e CacheError) String() string {
	return e.Message
}
