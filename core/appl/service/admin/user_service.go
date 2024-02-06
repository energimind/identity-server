package admin

import (
	"context"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
)

// UserService is a service for managing users.
//
// It implements the admin.UserService interface.
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

// Ensure service implements the admin.UserService interface.
var _ admin.UserService = (*UserService)(nil)

// GetUsers implements the admin.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) GetUsers(ctx context.Context, actor admin.Actor, appID admin.ID) ([]admin.User, error) {
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

// GetUser implements the admin.UserService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *UserService) GetUser(ctx context.Context, actor admin.Actor, appID, id admin.ID) (admin.User, error) {
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

// CreateUser implements the admin.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) CreateUser(ctx context.Context, actor admin.Actor, user admin.User) (admin.User, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.User{}, domain.NewAccessDeniedError("user %s cannot create user", actor.UserID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != user.ApplicationID {
			return admin.User{}, domain.NewAccessDeniedError("manager %s cannot create user", actor.UserID)
		}

		user.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleAdmin:
		user.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleNone:
		return admin.User{}, domain.NewAccessDeniedError("anonymous user cannot create user")
	default:
		return admin.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateUser implements the admin.UserService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *UserService) UpdateUser(ctx context.Context, actor admin.Actor, user admin.User) (admin.User, error) {
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

		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleAdmin:
		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return admin.User{}, err
		}

		return user, nil
	case admin.SystemRoleNone:
		return admin.User{}, domain.NewAccessDeniedError("anonymous user cannot update user %s", user.ID)
	default:
		return admin.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteUser implements the admin.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) DeleteUser(ctx context.Context, actor admin.Actor, appID, id admin.ID) error {
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

// GetUserByEmail implements the admin.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) GetUserByEmail(ctx context.Context, actor admin.Actor, appID admin.ID, email string) (admin.User, error) {
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
