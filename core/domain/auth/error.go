package auth

import "fmt"

// UserSessionError represents an error that occurs when handling user sessions.
type UserSessionError struct {
	Message string
}

// NewUserSessionError creates a new UserSessionError with the given message.
func NewUserSessionError(format string, args ...any) UserSessionError {
	return UserSessionError{Message: fmt.Sprintf(format, args...)}
}

// Error returns the error message.
func (e UserSessionError) Error() string {
	return fmt.Sprintf("User session error: %s", e.Message)
}

// String returns the error message.
func (e UserSessionError) String() string {
	return e.Error()
}
