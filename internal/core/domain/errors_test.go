package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	tester := func(t *testing.T, err, exp error) {
		require.ErrorAs(t, err, &exp)
		require.Equal(t, "test:42", err.Error())
	}

	tester(t, NewBadRequestError("test:%d", 42), BadRequestError{})
	tester(t, NewAccessDeniedError("test:%d", 42), AccessDeniedError{})
	tester(t, NewNotFoundError("test:%d", 42), NotFoundError{})
	tester(t, NewValidationError("test:%d", 42), ValidationError{})
	tester(t, NewStoreError("test:%d", 42), StoreError{})
	tester(t, NewConflictError("test:%d", 42), ConflictError{})
	tester(t, NewGatewayError("test:%d", 42), GatewayError{})
	tester(t, NewSessionError("test:%d", 42), SessionError{})
	tester(t, NewUnauthorizedError("test:%d", 42), UnauthorizedError{})
}
