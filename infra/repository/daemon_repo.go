package repository

import (
	"context"
	"errors"

	"github.com/energimind/identity-service/domain"
	"github.com/energimind/identity-service/domain/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DaemonRepository is a MongoDB implementation of DaemonRepository.
type DaemonRepository struct {
	db *mongo.Database
}

// NewDaemonRepository creates a new MongoDB daemon repository.
func NewDaemonRepository(db *mongo.Database) *DaemonRepository {
	return &DaemonRepository{db: db}
}

// Ensure repository implements the auth.DaemonRepository interface.
var _ auth.DaemonRepository = (*DaemonRepository)(nil)

// GetDaemons implements the auth.DaemonRepository interface.
func (r *DaemonRepository) GetDaemons(
	ctx context.Context,
	appID auth.ID,
) ([]auth.Daemon, error) {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"applicationId": appID}

	qCursor, err := coll.Find(ctx, qFilter)
	if err != nil {
		return nil, domain.NewStoreError("failed to find daemons: %v", err)
	}

	daemons, err := drainCursor[dbDaemon](ctx, qCursor, fromDaemon)
	if err != nil {
		return nil, domain.NewStoreError("failed to get daemons: %v", err)
	}

	return daemons, nil
}

// GetDaemon implements the auth.DaemonRepository interface.
func (r *DaemonRepository) GetDaemon(
	ctx context.Context,
	id auth.ID,
) (auth.Daemon, error) {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"id": id}
	daemon := dbDaemon{}

	if err := coll.FindOne(ctx, qFilter).Decode(&daemon); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return auth.Daemon{}, domain.NewNotFoundError("daemon %v not found", id)
		}

		return auth.Daemon{}, domain.NewStoreError("failed to get daemon: %v", err)
	}

	return fromDaemon(daemon), nil
}

// CreateDaemon implements the auth.DaemonRepository interface.
func (r *DaemonRepository) CreateDaemon(
	ctx context.Context,
	daemon auth.Daemon,
) error {
	coll := r.db.Collection("daemons")

	if _, err := coll.InsertOne(ctx, toDaemon(daemon)); err != nil {
		return domain.NewStoreError("failed to create daemon: %v", err)
	}

	return nil
}

// UpdateDaemon implements the auth.DaemonRepository interface.
func (r *DaemonRepository) UpdateDaemon(
	ctx context.Context,
	daemon auth.Daemon,
) error {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"id": daemon.ID}
	qUpdate := bson.M{"$set": toDaemon(daemon)}

	result, err := coll.UpdateOne(ctx, qFilter, qUpdate)
	if err != nil {
		return domain.NewStoreError("failed to update daemon: %v", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewNotFoundError("daemon %v not found", daemon.ID)
	}

	return nil
}

// DeleteDaemon implements the auth.DaemonRepository interface.
func (r *DaemonRepository) DeleteDaemon(
	ctx context.Context,
	id auth.ID,
) error {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"id": id}

	result, err := coll.DeleteOne(ctx, qFilter)
	if err != nil {
		return domain.NewStoreError("failed to delete daemon: %v", err)
	}

	if result.DeletedCount == 0 {
		return domain.NewNotFoundError("daemon %v not found", id)
	}

	return nil
}
