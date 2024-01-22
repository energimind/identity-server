package repository

import (
	"context"
	"errors"

	"github.com/energimind/identity-service/domain"
	"github.com/energimind/identity-service/domain/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProviderRepository is a MongoDB implementation of auth.ProviderRepository.
type ProviderRepository struct {
	db *mongo.Database
}

// NewProviderRepository creates a new MongoDB provider repository.
func NewProviderRepository(db *mongo.Database) *ProviderRepository {
	return &ProviderRepository{db: db}
}

// Ensure repository implements the auth.ProviderRepository interface.
var _ auth.ProviderRepository = (*ProviderRepository)(nil)

// GetProviders implements the auth.ProviderRepository interface.
func (r *ProviderRepository) GetProviders(
	ctx context.Context,
	appID auth.ID,
) ([]auth.Provider, error) {
	coll := r.db.Collection("providers")
	qFilter := bson.M{"applicationId": appID}

	qCursor, err := coll.Find(ctx, qFilter)
	if err != nil {
		return nil, domain.NewStoreError("failed to find providers: %v", err)
	}

	providers, err := drainCursor[dbProvider](ctx, qCursor, fromProvider)
	if err != nil {
		return nil, domain.NewStoreError("failed to get providers: %v", err)
	}

	return providers, nil
}

// GetProvider implements the auth.ProviderRepository interface.
func (r *ProviderRepository) GetProvider(
	ctx context.Context,
	id auth.ID,
) (auth.Provider, error) {
	coll := r.db.Collection("providers")
	qFilter := bson.M{"id": id}
	provider := dbProvider{}

	if err := coll.FindOne(ctx, qFilter).Decode(&provider); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return auth.Provider{}, domain.NewNotFoundError("provider %v not found", id)
		}

		return auth.Provider{}, domain.NewStoreError("failed to get provider: %v", err)
	}

	return fromProvider(provider), nil
}

// CreateProvider implements the auth.ProviderRepository interface.
func (r *ProviderRepository) CreateProvider(
	ctx context.Context,
	provider auth.Provider,
) error {
	coll := r.db.Collection("providers")

	if _, err := coll.InsertOne(ctx, toProvider(provider)); err != nil {
		return domain.NewStoreError("failed to create provider: %v", err)
	}

	return nil
}

// UpdateProvider implements the auth.ProviderRepository interface.
func (r *ProviderRepository) UpdateProvider(
	ctx context.Context,
	provider auth.Provider,
) error {
	coll := r.db.Collection("providers")
	qFilter := bson.M{"id": provider.ID}
	qUpdate := bson.M{"$set": toProvider(provider)}

	result, err := coll.UpdateOne(ctx, qFilter, qUpdate)
	if err != nil {
		return domain.NewStoreError("failed to update provider: %v", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewNotFoundError("provider %v not found", provider.ID)
	}

	return nil
}

// DeleteProvider implements the auth.ProviderRepository interface.
func (r *ProviderRepository) DeleteProvider(
	ctx context.Context,
	appID, id auth.ID,
) error {
	coll := r.db.Collection("providers")
	qFilter := bson.M{"id": id, "applicationId": appID}

	result, err := coll.DeleteOne(ctx, qFilter)
	if err != nil {
		return domain.NewStoreError("failed to delete provider: %v", err)
	}

	if result.DeletedCount == 0 {
		return domain.NewNotFoundError("provider %v not found", id)
	}

	return nil
}
