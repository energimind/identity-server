package service

import (
	"context"

	"github.com/energimind/identity-service/domain"
	"github.com/energimind/identity-service/domain/auth"
)

// DaemonService is a service for managing daemons.
//
// It implements the auth.DaemonService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
//
// Some methods are reported as to complex by the linter. We disable the linter for
// these methods, because they are not too complex, but just have a lot of error handling.
type DaemonService struct {
	repo  auth.DaemonRepository
	idgen domain.IDGenerator
}

// NewDaemonService returns a new DaemonService instance.
func NewDaemonService(
	repo auth.DaemonRepository,
	idgen domain.IDGenerator,
) *DaemonService {
	return &DaemonService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the auth.DaemonService interface.
var _ auth.DaemonService = (*DaemonService)(nil)

// GetDaemons implements the auth.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) GetDaemons(
	ctx context.Context,
	actor auth.Actor,
	appID auth.ID,
) ([]auth.Daemon, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get daemons", actor.UserID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != appID {
			return nil, domain.NewAccessDeniedError("manager %s cannot get daemons for application %s", actor.UserID, appID)
		}

		daemons, err := s.repo.GetDaemons(ctx, appID)
		if err != nil {
			return nil, err
		}

		return daemons, nil
	case auth.SystemRoleAdmin:
		daemons, err := s.repo.GetDaemons(ctx, appID)
		if err != nil {
			return nil, err
		}

		return daemons, nil
	case auth.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get daemons")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetDaemon implements the auth.DaemonService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *DaemonService) GetDaemon(
	ctx context.Context,
	actor auth.Actor,
	appID, id auth.ID,
) (auth.Daemon, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Daemon{}, domain.NewAccessDeniedError("user %s cannot get daemon %s", actor.UserID, id)
	case auth.SystemRoleManager:
		daemon, err := s.repo.GetDaemon(ctx, id)
		if err != nil {
			return auth.Daemon{}, err
		}

		if actor.ApplicationID != appID {
			return auth.Daemon{}, domain.NewAccessDeniedError("manager %s cannot get daemon %s", actor.UserID, id)
		}

		return daemon, nil
	case auth.SystemRoleAdmin:
		user, err := s.repo.GetDaemon(ctx, id)
		if err != nil {
			return auth.Daemon{}, err
		}

		return user, nil
	case auth.SystemRoleNone:
		return auth.Daemon{}, domain.NewAccessDeniedError("anonymous user cannot get daemon %s", id)
	default:
		return auth.Daemon{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateDaemon implements the auth.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) CreateDaemon(
	ctx context.Context,
	actor auth.Actor,
	daemon auth.Daemon,
) (auth.Daemon, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Daemon{}, domain.NewAccessDeniedError("user %s cannot create daemon", actor.UserID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != daemon.ApplicationID {
			return auth.Daemon{}, domain.NewAccessDeniedError("manager %s cannot create daemon", actor.UserID)
		}

		daemon.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateDaemon(ctx, daemon); err != nil {
			return auth.Daemon{}, err
		}

		return daemon, nil
	case auth.SystemRoleAdmin:
		daemon.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateDaemon(ctx, daemon); err != nil {
			return auth.Daemon{}, err
		}

		return daemon, nil
	case auth.SystemRoleNone:
		return auth.Daemon{}, domain.NewAccessDeniedError("anonymous user cannot create daemon")
	default:
		return auth.Daemon{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateDaemon implements the auth.DaemonService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *DaemonService) UpdateDaemon(
	ctx context.Context,
	actor auth.Actor,
	daemon auth.Daemon,
) (auth.Daemon, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Daemon{}, domain.NewAccessDeniedError("user %s cannot update daemon %s", actor.UserID, daemon.ID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != daemon.ApplicationID {
			return auth.Daemon{}, domain.NewAccessDeniedError("manager %s cannot update daemon %s", actor.UserID, daemon.ID)
		}

		if err := s.repo.UpdateDaemon(ctx, daemon); err != nil {
			return auth.Daemon{}, err
		}

		return daemon, nil
	case auth.SystemRoleAdmin:
		if err := s.repo.UpdateDaemon(ctx, daemon); err != nil {
			return auth.Daemon{}, err
		}

		return daemon, nil
	case auth.SystemRoleNone:
		return auth.Daemon{}, domain.NewAccessDeniedError("anonymous user cannot update daemon %s", daemon.ID)
	default:
		return auth.Daemon{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteDaemon implements the auth.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) DeleteDaemon(
	ctx context.Context,
	actor auth.Actor,
	appID, id auth.ID,
) error {
	switch actor.Role {
	case auth.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete daemon %s", actor.UserID, id)
	case auth.SystemRoleManager:
		if actor.ApplicationID != appID {
			return domain.NewAccessDeniedError("manager %s cannot delete daemon %s", actor.UserID, id)
		}

		if err := s.repo.DeleteDaemon(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case auth.SystemRoleAdmin:
		if err := s.repo.DeleteDaemon(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case auth.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete daemon %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
