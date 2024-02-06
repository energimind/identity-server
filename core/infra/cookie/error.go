package cookie

import "fmt"

// Error represents an error that occurs when handling cookies.
type Error struct {
	Message string
}

// NewError creates a new Error with the given message.
func NewError(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("Cookie error: %s", e.Message)
}
