package service

import (
	"context"

	"github.com/energimind/identity-service/domain"
	"github.com/energimind/identity-service/domain/auth"
)

// UserService is a service for managing users.
//
// It implements the auth.UserService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
//
// Some methods are reported as to complex by the linter. We disable the linter for
// these methods, because they are not too complex, but just have a lot of error handling.
type UserService struct {
	repo  auth.UserRepository
	idgen domain.IDGenerator
}

// NewUserService returns a new UserService instance.
func NewUserService(
	repo auth.UserRepository,
	idgen domain.IDGenerator,
) *UserService {
	return &UserService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the auth.UserService interface.
var _ auth.UserService = (*UserService)(nil)

// GetUsers implements the auth.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) GetUsers(ctx context.Context, actor auth.Actor, appID auth.ID) ([]auth.User, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get users", actor.UserID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != appID {
			return nil, domain.NewAccessDeniedError("manager %s cannot get users for application %s", actor.UserID, appID)
		}

		users, err := s.repo.GetUsers(ctx, appID)
		if err != nil {
			return nil, err
		}

		return users, nil
	case auth.SystemRoleAdmin:
		users, err := s.repo.GetUsers(ctx, appID)
		if err != nil {
			return nil, err
		}

		return users, nil
	case auth.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get users")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetUser implements the auth.UserService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *UserService) GetUser(ctx context.Context, actor auth.Actor, appID, id auth.ID) (auth.User, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		if actor.ApplicationID != appID {
			return auth.User{}, domain.NewAccessDeniedError("user %s cannot get user %s", actor.UserID, id)
		}

		// user can only get itself
		if actor.UserID != id {
			return auth.User{}, domain.NewAccessDeniedError("user %s cannot get user %s", actor.UserID, id)
		}

		user, err := s.repo.GetUser(ctx, id)
		if err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleManager:
		user, err := s.repo.GetUser(ctx, id)
		if err != nil {
			return auth.User{}, err
		}

		if actor.ApplicationID != appID {
			return auth.User{}, domain.NewAccessDeniedError("manager %s cannot get user %s", actor.UserID, id)
		}

		return user, nil
	case auth.SystemRoleAdmin:
		user, err := s.repo.GetUser(ctx, id)
		if err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleNone:
		return auth.User{}, domain.NewAccessDeniedError("anonymous user cannot get user %s", id)
	default:
		return auth.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateUser implements the auth.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) CreateUser(ctx context.Context, actor auth.Actor, user auth.User) (auth.User, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.User{}, domain.NewAccessDeniedError("user %s cannot create user", actor.UserID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != user.ApplicationID {
			return auth.User{}, domain.NewAccessDeniedError("manager %s cannot create user", actor.UserID)
		}

		user.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleAdmin:
		user.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateUser(ctx, user); err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleNone:
		return auth.User{}, domain.NewAccessDeniedError("anonymous user cannot create user")
	default:
		return auth.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateUser implements the auth.UserService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *UserService) UpdateUser(ctx context.Context, actor auth.Actor, user auth.User) (auth.User, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		if actor.ApplicationID != user.ApplicationID {
			return auth.User{}, domain.NewAccessDeniedError("user %s cannot update user %s", actor.UserID, user.ID)
		}

		// user can only update itself
		if actor.UserID != user.ID {
			return auth.User{}, domain.NewAccessDeniedError("user %s cannot update user %s", actor.UserID, user.ID)
		}

		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleManager:
		if actor.ApplicationID != user.ApplicationID {
			return auth.User{}, domain.NewAccessDeniedError("manager %s cannot update user %s", actor.UserID, user.ID)
		}

		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleAdmin:
		if err := s.repo.UpdateUser(ctx, user); err != nil {
			return auth.User{}, err
		}

		return user, nil
	case auth.SystemRoleNone:
		return auth.User{}, domain.NewAccessDeniedError("anonymous user cannot update user %s", user.ID)
	default:
		return auth.User{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteUser implements the auth.UserService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *UserService) DeleteUser(ctx context.Context, actor auth.Actor, appID, id auth.ID) error {
	switch actor.Role {
	case auth.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete user %s", actor.UserID, id)
	case auth.SystemRoleManager:
		if actor.ApplicationID != appID {
			return domain.NewAccessDeniedError("manager %s cannot delete user %s", actor.UserID, id)
		}

		if err := s.repo.DeleteUser(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case auth.SystemRoleAdmin:
		if err := s.repo.DeleteUser(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case auth.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete user %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
