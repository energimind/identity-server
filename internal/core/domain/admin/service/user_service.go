package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// UserService is a service for managing users.
//
// It implements the service.UserService and the admin.UserFinder interfaces.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
//
// Some methods are reported as to complex by the linter. We disable the linter for
// these methods, because they are not too complex, but just have a lot of error handling.
type UserService struct {
	repo  admin.UserRepository
	idgen domain.IDGenerator
}

// NewUserService returns a new UserService instance.
func NewUserService(
	repo admin.UserRepository,
	idgen domain.IDGenerator,
) *UserService {
	return &UserService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the service.UserService interface.
var _ admin.UserService = (*UserService)(nil)

// Ensure service implements the admin.UserFinder interface.
var _ admin.UserFinder = (*UserService)(nil)

// GetUsers implements the service.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) GetUsers(
	ctx context.Context,
	actor admin.Actor,
	appID admin.ID,
) ([]admin.User, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get users", actor.UserID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return nil, domain.NewAccessDeniedError("manager %s cannot get users for application %s", actor.UserID, appID)
		}

		users, err := s.repo.GetUsers(ctx, appID)
		if err != nil {
			return nil, err
		}

		return users, nil
	case admin.SystemRoleAdmin:
		users, err := s.repo.GetUsers(ctx, appID)
		if err != nil {
			return nil, err
		}

		return users, nil
	case admin.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get users")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetUser implements the service.UserService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *UserService) GetUser(
	ctx context.Context,
	actor admin.Actor,
	appID, id admin.ID,
) (admin.User, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		if actor.ApplicationID != appID {
			return admin.User{}, domain.NewAccessDeniedError("user %s cannot get user %s", actor.UserID, id)
		}

		// user can only get itself
		if actor.UserID != id {
			return admin.User{}, domain.NewAccessDeniedError("user %s cannot get user %s", actor.UserID, id)
		}

		user, err := s.repo.GetUser(ctx, appID, id)
		if err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return admin.User{}, domain.NewAccessDeniedError("manager %s cannot get user %s", actor.UserID, id)
		}

		user, err := s.repo.GetUser(ctx, appID, id)
		if err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleAdmin:
		user, err := s.repo.GetUser(ctx, appID, id)
		if err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleNone:
		return admin.User{}, domain.NewAccessDeniedError("anonymous user cannot get user %s", id)
	default:
		return admin.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateUser implements the service.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) CreateUser(
	ctx context.Context,
	actor admin.Actor,
	user admin.User,
) (admin.User, error) {
	user, err := validateUser(user)
	if err != nil {
		return admin.User{}, err
	}

	create := func() (admin.User, error) {
		if err := s.checkUserExists(ctx, user.ApplicationID, user.Email); err != nil {
			return admin.User{}, err
		}

		user.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.User{}, domain.NewAccessDeniedError("user %s cannot create user", actor.UserID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != user.ApplicationID {
			return admin.User{}, domain.NewAccessDeniedError("manager %s cannot create user", actor.UserID)
		}

		return create()
	case admin.SystemRoleAdmin:
		return create()
	case admin.SystemRoleNone:
		return admin.User{}, domain.NewAccessDeniedError("anonymous user cannot create user")
	default:
		return admin.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateUser implements the service.UserService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *UserService) UpdateUser(
	ctx context.Context,
	actor admin.Actor,
	user admin.User,
) (admin.User, error) {
	user, err := validateUser(user)
	if err != nil {
		return admin.User{}, err
	}

	update := func() (admin.User, error) {
		if err := s.checkAnotherUserExists(ctx, user.ApplicationID, user.Email, user.ID); err != nil {
			return admin.User{}, err
		}

		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		if actor.ApplicationID != user.ApplicationID {
			return admin.User{}, domain.NewAccessDeniedError("user %s cannot update user %s", actor.UserID, user.ID)
		}

		// user can only update itself
		if actor.UserID != user.ID {
			return admin.User{}, domain.NewAccessDeniedError("user %s cannot update user %s", actor.UserID, user.ID)
		}

		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleManager:
		if actor.ApplicationID != user.ApplicationID {
			return admin.User{}, domain.NewAccessDeniedError("manager %s cannot update user %s", actor.UserID, user.ID)
		}

		return update()
	case admin.SystemRoleAdmin:
		return update()
	case admin.SystemRoleNone:
		return admin.User{}, domain.NewAccessDeniedError("anonymous user cannot update user %s", user.ID)
	default:
		return admin.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteUser implements the service.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) DeleteUser(
	ctx context.Context,
	actor admin.Actor,
	appID, id admin.ID,
) error {
	switch actor.Role {
	case admin.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete user %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return domain.NewAccessDeniedError("manager %s cannot delete user %s", actor.UserID, id)
		}

		if err := s.repo.DeleteUser(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleAdmin:
		if err := s.repo.DeleteUser(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete user %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetAPIKeys implements the service.UserService interface.
func (s *UserService) GetAPIKeys(
	ctx context.Context,
	actor admin.Actor,
	appID, userID admin.ID,
) ([]admin.APIKey, error) {
	user, err := s.GetUser(ctx, actor, appID, userID)
	if err != nil {
		return nil, err
	}

	return user.APIKeys, nil
}

// GetAPIKey implements the service.UserService interface.
func (s *UserService) GetAPIKey(
	ctx context.Context,
	actor admin.Actor,
	appID, userID, id admin.ID,
) (admin.APIKey, error) {
	user, err := s.GetUser(ctx, actor, appID, userID)
	if err != nil {
		return admin.APIKey{}, err
	}

	for _, apiKey := range user.APIKeys {
		if apiKey.ID == id {
			return apiKey, nil
		}
	}

	return admin.APIKey{}, domain.NewNotFoundError("API key %s not found", id)
}

// CreateAPIKey implements the service.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) CreateAPIKey(
	ctx context.Context,
	actor admin.Actor,
	appID, userID admin.ID,
	apiKey admin.APIKey,
) (admin.APIKey, error) {
	apiKey, err := validateAPIKey(apiKey)
	if err != nil {
		return admin.APIKey{}, err
	}

	user, err := s.GetUser(ctx, actor, appID, userID)
	if err != nil {
		return admin.APIKey{}, err
	}

	apiKey.ID = admin.ID(s.idgen.GenerateID())

	user.APIKeys = append(user.APIKeys, apiKey)

	if uErr := s.repo.UpdateUser(ctx, user); uErr != nil {
		return admin.APIKey{}, uErr
	}

	return apiKey, nil
}

// UpdateAPIKey implements the service.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) UpdateAPIKey(
	ctx context.Context,
	actor admin.Actor,
	appID, userID, id admin.ID,
	apiKey admin.APIKey,
) (admin.APIKey, error) {
	apiKey, err := validateAPIKey(apiKey)
	if err != nil {
		return admin.APIKey{}, err
	}

	user, err := s.GetUser(ctx, actor, appID, userID)
	if err != nil {
		return admin.APIKey{}, err
	}

	for i, ak := range user.APIKeys {
		if ak.ID == id {
			user.APIKeys[i] = apiKey

			if uErr := s.repo.UpdateUser(ctx, user); uErr != nil {
				return admin.APIKey{}, uErr
			}

			return apiKey, nil
		}
	}

	return admin.APIKey{}, domain.NewNotFoundError("API key %s not found", id)
}

// DeleteAPIKey implements the service.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) DeleteAPIKey(
	ctx context.Context,
	actor admin.Actor,
	appID, userID, id admin.ID,
) error {
	user, err := s.GetUser(ctx, actor, appID, userID)
	if err != nil {
		return err
	}

	for i, apiKey := range user.APIKeys {
		if apiKey.ID == id {
			user.APIKeys = append(user.APIKeys[:i], user.APIKeys[i+1:]...)

			if uErr := s.repo.UpdateUser(ctx, user); uErr != nil {
				return uErr
			}

			return nil
		}
	}

	return domain.NewNotFoundError("API key %s not found", id)
}

// GetUserByEmail implements the admin.UserFinder interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) GetUserByEmail(
	ctx context.Context,
	actor admin.Actor,
	appID admin.ID,
	email string,
) (admin.User, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.User{}, domain.NewAccessDeniedError("user %s cannot get user by email", actor.UserID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return admin.User{}, domain.NewAccessDeniedError("manager %s cannot get user by email", actor.UserID)
		}

		user, err := s.repo.GetUserByEmail(ctx, appID, email)
		if err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleAdmin:
		user, err := s.repo.GetUserByEmail(ctx, appID, email)
		if err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleNone:
		return admin.User{}, domain.NewAccessDeniedError("anonymous user cannot get user by email")
	default:
		return admin.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// checkUserExists checks if a user with the given email already exists.
//
// It returns a domain.ConflictError if the user already exists.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) checkUserExists(ctx context.Context, appID admin.ID, email string) error {
	_, err := s.repo.GetUserByEmail(ctx, appID, email)
	if err == nil {
		return domain.NewConflictError("user with email %s already exists", email)
	}

	if domain.IsNotFoundError(err) {
		return nil
	}

	return err
}

// checkAnotherUserExists checks if a user with the given email already exists, but not the user with the given ID.
//
// It returns a domain.ConflictError if the user already exists.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) checkAnotherUserExists(ctx context.Context, appID admin.ID, email string, id admin.ID) error {
	user, err := s.repo.GetUserByEmail(ctx, appID, email)
	if err != nil {
		return err
	}

	if user.ID != id {
		return domain.NewConflictError("user with email %s already exists", email)
	}

	return nil
}
