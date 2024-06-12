package repository

import (
	"context"
	"errors"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
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

// Ensure repository implements the admin.DaemonRepository interface.
var _ admin.DaemonRepository = (*DaemonRepository)(nil)

// GetDaemons implements the admin.DaemonRepository interface.
func (r *DaemonRepository) GetDaemons(
	ctx context.Context,
	realmID admin.ID,
) ([]admin.Daemon, error) {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"realmId": realmID}

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

// GetDaemon implements the admin.DaemonRepository interface.
func (r *DaemonRepository) GetDaemon(
	ctx context.Context,
	realmID, id admin.ID,
) (admin.Daemon, error) {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"id": id, "realmId": realmID}
	daemon := dbDaemon{}

	if err := coll.FindOne(ctx, qFilter).Decode(&daemon); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return admin.Daemon{}, domain.NewNotFoundError("daemon %v not found", id)
		}

		return admin.Daemon{}, domain.NewStoreError("failed to get daemon: %v", err)
	}

	return fromDaemon(daemon), nil
}

// CreateDaemon implements the admin.DaemonRepository interface.
func (r *DaemonRepository) CreateDaemon(
	ctx context.Context,
	daemon admin.Daemon,
) error {
	coll := r.db.Collection("daemons")

	if _, err := coll.InsertOne(ctx, toDaemon(daemon)); err != nil {
		return domain.NewStoreError("failed to create daemon: %v", err)
	}

	return nil
}

// UpdateDaemon implements the admin.DaemonRepository interface.
func (r *DaemonRepository) UpdateDaemon(
	ctx context.Context,
	daemon admin.Daemon,
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

// DeleteDaemon implements the admin.DaemonRepository interface.
func (r *DaemonRepository) DeleteDaemon(
	ctx context.Context,
	realmID, id admin.ID,
) error {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{"id": id, "realmId": realmID}

	result, err := coll.DeleteOne(ctx, qFilter)
	if err != nil {
		return domain.NewStoreError("failed to delete daemon: %v", err)
	}

	if result.DeletedCount == 0 {
		return domain.NewNotFoundError("daemon %v not found", id)
	}

	return nil
}

// GetAPIKey implements the admin.UserRepository interface.
//
// This method takes in account the enabled field of the user and the API key.
func (r *DaemonRepository) GetAPIKey(
	ctx context.Context,
	realmID admin.ID,
	key string,
) (admin.APIKey, error) {
	coll := r.db.Collection("daemons")
	qFilter := bson.M{
		"realmId": realmID,
		"enabled": true,
		"apiKeys": bson.M{"$elemMatch": bson.M{
			"key":     key,
			"enabled": true,
		}},
	}
	daemon := dbDaemon{}

	if err := coll.FindOne(ctx, qFilter).Decode(&daemon); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return admin.APIKey{}, domain.NewNotFoundError("API key not found")
		}

		return admin.APIKey{}, domain.NewStoreError("failed to API api key: %v", err)
	}

	for _, apiKey := range daemon.APIKeys {
		if apiKey.Key == key {
			return fromAPIKey(apiKey), nil
		}
	}

	return admin.APIKey{}, domain.NewNotFoundError("API key not found")
}
