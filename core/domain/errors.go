package domain

import "fmt"

// BadRequestError is the error returned when a request is invalid.
type BadRequestError struct {
	Message string
}

// NewBadRequestError returns a new BadRequestError.
func NewBadRequestError(format string, args ...any) BadRequestError {
	return BadRequestError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e BadRequestError) Error() string {
	return e.Message
}

// AccessDeniedError is the error returned when a user is not allowed to perform an action.
type AccessDeniedError struct {
	Message string
}

// NewAccessDeniedError returns a new AccessDeniedError.
func NewAccessDeniedError(format string, args ...any) AccessDeniedError {
	return AccessDeniedError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e AccessDeniedError) Error() string {
	return e.Message
}

// NotFoundError is the error returned when an object is not found.
type NotFoundError struct {
	Message string
}

// NewNotFoundError returns a new NotFoundError.
func NewNotFoundError(format string, args ...any) NotFoundError {
	return NotFoundError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e NotFoundError) Error() string {
	return e.Message
}

// ValidationError is the error returned when an object is invalid.
type ValidationError struct {
	Message string
}

// NewValidationError returns a new ValidationError.
func NewValidationError(format string, args ...any) ValidationError {
	return ValidationError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e ValidationError) Error() string {
	return e.Message
}

// StoreError is the error returned when an error occurs while storing an object.
type StoreError struct {
	Message string
}

// NewStoreError returns a new StoreError.
func NewStoreError(format string, args ...any) StoreError {
	return StoreError{
		Message: fmt.Sprintf(format, args...),
	}
}

// Error returns the error message.
func (e StoreError) Error() string {
	return e.Message
}
