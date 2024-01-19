package repository

import (
	"context"

	"github.com/energimind/identity-service/domain/auth"
	"go.mongodb.org/mongo-driver/bson"
)

// filter is a helper type for building MongoDB filters.
type filter bson.M

// newFilter creates a new filter.
func newFilter() filter {
	return filter{}
}

// id adds a primary field to the filter.
func (f filter) id(id auth.ID) filter {
	f["id"] = id

	return f
}

// matchID adds a primary field to the filter if the principal ID matches the ID.
// If the principal ID is empty, the ID field is set to its original value.
// If the principal ID does not match the ID, the ID field is set to "0" to
// prevent the query from returning any results.
func (f filter) matchID(id, principalID auth.ID) filter {
	const denyAccess = "0"

	if principalID != "" && id != principalID {
		id = denyAccess
	}

	return f.id(id)
}

// scope adds a field to the filter if the value is not empty.
func (f filter) scope(field string, value auth.ID) filter {
	if !value.IsEmpty() {
		f[field] = value.String()
	}

	return f
}

// cursor is an interface for a MongoDB cursor.
// It is used to mock the mongo.Cursor type.
type cursor interface {
	All(ctx context.Context, v any) error
	Close(ctx context.Context) error
}

// drainCursor drains the cursor into a slice of T and then maps it to a slice of M.
//
//nolint:wrapcheck
func drainCursor[T, M any](ctx context.Context, cursor cursor, mapper func([]T) []M) ([]M, error) {
	const preallocate = 4

	results := make([]T, 0, preallocate)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}

	return mapper(results), nil
}
