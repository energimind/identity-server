package repository

import (
	"context"
	"errors"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RealmRepository is a MongoDB implementation of the admin.RealmRepository interface.
type RealmRepository struct {
	db *mongo.Database
}

// NewRealmRepository creates a new MongoDB realm repository.
func NewRealmRepository(db *mongo.Database) *RealmRepository {
	return &RealmRepository{db: db}
}

// Ensure repository implements the admin.RealmRepository interface.
var _ admin.RealmRepository = (*RealmRepository)(nil)

// GetRealms implements the admin.RealmRepository interface.
func (r *RealmRepository) GetRealms(
	ctx context.Context,
) ([]admin.Realm, error) {
	coll := r.db.Collection("realms")
	qFilter := bson.M{}

	qCursor, err := coll.Find(ctx, qFilter)
	if err != nil {
		return nil, domain.NewStoreError("failed to find realms: %v", err)
	}

	realms, err := drainCursor[dbRealm](ctx, qCursor, fromRealm)
	if err != nil {
		return nil, domain.NewStoreError("failed to get realms: %v", err)
	}

	return realms, nil
}

// GetRealm implements the admin.RealmRepository interface.
func (r *RealmRepository) GetRealm(
	ctx context.Context,
	id admin.ID,
) (admin.Realm, error) {
	coll := r.db.Collection("realms")
	qFilter := bson.M{"id": id}
	realm := dbRealm{}

	if err := coll.FindOne(ctx, qFilter).Decode(&realm); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return admin.Realm{}, domain.NewNotFoundError("realm %v not found", id)
		}

		return admin.Realm{}, domain.NewStoreError("failed to get realm: %v", err)
	}

	return fromRealm(realm), nil
}

// CreateRealm implements the admin.RealmRepository interface.
func (r *RealmRepository) CreateRealm(
	ctx context.Context,
	realm admin.Realm,
) error {
	coll := r.db.Collection("realms")

	if _, err := coll.InsertOne(ctx, toRealm(realm)); err != nil {
		return domain.NewStoreError("failed to create realm: %v", err)
	}

	return nil
}

// UpdateRealm implements the admin.RealmRepository interface.
func (r *RealmRepository) UpdateRealm(
	ctx context.Context,
	realm admin.Realm,
) error {
	coll := r.db.Collection("realms")
	qFilter := bson.M{"id": realm.ID}
	qUpdate := bson.M{"$set": toRealm(realm)}

	result, err := coll.UpdateOne(ctx, qFilter, qUpdate)
	if err != nil {
		return domain.NewStoreError("failed to update realm: %v", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewNotFoundError("realm %v not found", realm.ID)
	}

	return nil
}

// DeleteRealm implements the admin.RealmRepository interface.
func (r *RealmRepository) DeleteRealm(
	ctx context.Context,
	id admin.ID,
) error {
	coll := r.db.Collection("realms")
	qFilter := bson.M{"id": id}

	result, err := coll.DeleteOne(ctx, qFilter)
	if err != nil {
		return domain.NewStoreError("failed to delete realm: %v", err)
	}

	if result.DeletedCount == 0 {
		return domain.NewNotFoundError("realm %v not found", id)
	}

	return nil
}
