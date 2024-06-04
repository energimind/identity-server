package client

import "fmt"

// IdentityServerError is the error returned when an error occurs while communicating with
// the identity service.
type IdentityServerError struct {
	Message string
}

// newIdentityServerError returns a new IdentityServerError.
func newIdentityServerError(format string, args ...any) IdentityServerError {
	return IdentityServerError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e IdentityServerError) Error() string {
	return e.Message
}
