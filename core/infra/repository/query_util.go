package repository

import (
	"context"
)

// cursor is an interface for a MongoDB cursor.
// It is used to mock the mongo.Cursor type.
type cursor interface {
	All(ctx context.Context, v any) error
	Close(ctx context.Context) error
}

// drainCursor drains the cursor into a slice of T and then maps it to a slice of M.
//
//nolint:wrapcheck
func drainCursor[T, M any](ctx context.Context, cursor cursor, mapper func(T) M) ([]M, error) {
	const preallocate = 4

	results := make([]T, 0, preallocate)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}

	mapped := make([]M, len(results))

	for i, result := range results {
		mapped[i] = mapper(result)
	}

	return mapped, nil
}
