package server

import (
	"context"
	"time"

	"github.com/energimind/go-kit/slog"
	"github.com/pkg/errors"
)

// withDelayReporter reports when an operation takes longer than expected.
//
//nolint:ireturn
func withDelayReporter[T any](
	ctx context.Context,
	fn func() (T, error),
	operation string,
) (T, error) {
	const timeout = 2 * time.Second

	var result struct {
		value T
		err   error
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		result.value, result.err = fn()
	}()

	messageDisplayed := false

	for {
		select {
		case <-ctx.Done():
			return result.value, errors.Wrap(result.err, "context cancelled")
		case <-done:
			return result.value, result.err
		case <-time.After(timeout):
			if !messageDisplayed {
				messageDisplayed = true

				slog.Warn().Msgf("Operation %s is taking longer than expected", operation)
			}
		}
	}
}
