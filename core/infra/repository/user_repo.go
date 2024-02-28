package repository

import (
	"context"
	"errors"

	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository is a MongoDB implementation of UserRepository.
type UserRepository struct {
	db *mongo.Database
}

// NewUserRepository creates a new MongoDB user repository.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

// Ensure repository implements the admin.UserRepository interface.
var _ admin.UserRepository = (*UserRepository)(nil)

// GetUsers implements the admin.UserRepository interface.
func (r *UserRepository) GetUsers(
	ctx context.Context,
	appID admin.ID,
) ([]admin.User, error) {
	coll := r.db.Collection("users")
	qFilter := bson.M{"applicationId": appID}

	qCursor, err := coll.Find(ctx, qFilter)
	if err != nil {
		return nil, domain.NewStoreError("failed to find users: %v", err)
	}

	users, err := drainCursor[dbUser](ctx, qCursor, fromUser)
	if err != nil {
		return nil, domain.NewStoreError("failed to get users: %v", err)
	}

	return users, nil
}

// GetUser implements the admin.UserRepository interface.
func (r *UserRepository) GetUser(
	ctx context.Context,
	appID, id admin.ID,
) (admin.User, error) {
	coll := r.db.Collection("users")
	qFilter := bson.M{"id": id, "applicationId": appID}
	user := dbUser{}

	if err := coll.FindOne(ctx, qFilter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return admin.User{}, domain.NewNotFoundError("user %v not found", id)
		}

		return admin.User{}, domain.NewStoreError("failed to get user: %v", err)
	}

	return fromUser(user), nil
}

// CreateUser implements the admin.UserRepository interface.
func (r *UserRepository) CreateUser(
	ctx context.Context,
	user admin.User,
) error {
	coll := r.db.Collection("users")

	if _, err := coll.InsertOne(ctx, toUser(user)); err != nil {
		return domain.NewStoreError("failed to create user: %v", err)
	}

	return nil
}

// UpdateUser implements the admin.UserRepository interface.
func (r *UserRepository) UpdateUser(
	ctx context.Context,
	user admin.User,
) error {
	coll := r.db.Collection("users")
	qFilter := bson.M{"id": user.ID}
	qUpdate := bson.M{"$set": toUser(user)}

	result, err := coll.UpdateOne(ctx, qFilter, qUpdate)
	if err != nil {
		return domain.NewStoreError("failed to update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewNotFoundError("user %v not found", user.ID)
	}

	return nil
}

// DeleteUser implements the admin.UserRepository interface.
func (r *UserRepository) DeleteUser(
	ctx context.Context,
	appID, id admin.ID,
) error {
	coll := r.db.Collection("users")
	qFilter := bson.M{"id": id, "applicationId": appID}

	result, err := coll.DeleteOne(ctx, qFilter)
	if err != nil {
		return domain.NewStoreError("failed to delete user: %v", err)
	}

	if result.DeletedCount == 0 {
		return domain.NewNotFoundError("user %v not found", id)
	}

	return nil
}

// GetUserByEmail implements the admin.UserRepository interface.
func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	appID admin.ID,
	email string,
) (admin.User, error) {
	coll := r.db.Collection("users")
	qFilter := bson.M{"email": email, "applicationId": appID}
	user := dbUser{}

	if err := coll.FindOne(ctx, qFilter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return admin.User{}, domain.NewNotFoundError("user with email %s not found", email)
		}

		return admin.User{}, domain.NewStoreError("failed to get user by email: %v", err)
	}

	return fromUser(user), nil
}
