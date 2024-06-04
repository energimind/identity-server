package oauth

import "fmt"

// Error represents an error returned by the OAuth package.
type Error struct {
	Message string
}

// NewError creates a new Error.
func NewError(format string, args ...any) *Error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}
