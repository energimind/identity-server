package sessioncookie

import "fmt"

// Error contains an error message.
type Error struct {
	Message string
}

// NewError returns a new instance of Error.
func NewError(format string, args ...any) Error {
	return Error{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e Error) Error() string {
	return e.Message
}
