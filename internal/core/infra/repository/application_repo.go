package repository

import (
	"context"
	"errors"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplicationRepository is a MongoDB implementation of the admin.ApplicationRepository interface.
type ApplicationRepository struct {
	db *mongo.Database
}

// NewApplicationRepository creates a new MongoDB application repository.
func NewApplicationRepository(db *mongo.Database) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// Ensure repository implements the admin.ApplicationRepository interface.
var _ admin.ApplicationRepository = (*ApplicationRepository)(nil)

// GetApplications implements the admin.ApplicationRepository interface.
func (r *ApplicationRepository) GetApplications(
	ctx context.Context,
) ([]admin.Application, error) {
	coll := r.db.Collection("applications")
	qFilter := bson.M{}

	qCursor, err := coll.Find(ctx, qFilter)
	if err != nil {
		return nil, domain.NewStoreError("failed to find applications: %v", err)
	}

	applications, err := drainCursor[dbApplication](ctx, qCursor, fromApplication)
	if err != nil {
		return nil, domain.NewStoreError("failed to get applications: %v", err)
	}

	return applications, nil
}

// GetApplication implements the admin.ApplicationRepository interface.
func (r *ApplicationRepository) GetApplication(
	ctx context.Context,
	id admin.ID,
) (admin.Application, error) {
	coll := r.db.Collection("applications")
	qFilter := bson.M{"id": id}
	application := dbApplication{}

	if err := coll.FindOne(ctx, qFilter).Decode(&application); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return admin.Application{}, domain.NewNotFoundError("application %v not found", id)
		}

		return admin.Application{}, domain.NewStoreError("failed to get application: %v", err)
	}

	return fromApplication(application), nil
}

// CreateApplication implements the admin.ApplicationRepository interface.
func (r *ApplicationRepository) CreateApplication(
	ctx context.Context,
	app admin.Application,
) error {
	coll := r.db.Collection("applications")

	if _, err := coll.InsertOne(ctx, toApplication(app)); err != nil {
		return domain.NewStoreError("failed to create application: %v", err)
	}

	return nil
}

// UpdateApplication implements the admin.ApplicationRepository interface.
func (r *ApplicationRepository) UpdateApplication(
	ctx context.Context,
	app admin.Application,
) error {
	coll := r.db.Collection("applications")
	qFilter := bson.M{"id": app.ID}
	qUpdate := bson.M{"$set": toApplication(app)}

	result, err := coll.UpdateOne(ctx, qFilter, qUpdate)
	if err != nil {
		return domain.NewStoreError("failed to update application: %v", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewNotFoundError("application %v not found", app.ID)
	}

	return nil
}

// DeleteApplication implements the admin.ApplicationRepository interface.
func (r *ApplicationRepository) DeleteApplication(
	ctx context.Context,
	id admin.ID,
) error {
	coll := r.db.Collection("applications")
	qFilter := bson.M{"id": id}

	result, err := coll.DeleteOne(ctx, qFilter)
	if err != nil {
		return domain.NewStoreError("failed to delete application: %v", err)
	}

	if result.DeletedCount == 0 {
		return domain.NewNotFoundError("application %v not found", id)
	}

	return nil
}
