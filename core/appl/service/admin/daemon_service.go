package admin

import (
	"context"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
)

// DaemonService is a service for managing daemons.
//
// It implements the admin.DaemonService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type DaemonService struct {
	repo  admin.DaemonRepository
	idgen domain.IDGenerator
}

// NewDaemonService returns a new DaemonService instance.
func NewDaemonService(
	repo admin.DaemonRepository,
	idgen domain.IDGenerator,
) *DaemonService {
	return &DaemonService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the admin.DaemonService interface.
var _ admin.DaemonService = (*DaemonService)(nil)

// GetDaemons implements the admin.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) GetDaemons(
	ctx context.Context,
	actor admin.Actor,
	appID admin.ID,
) ([]admin.Daemon, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get daemons", actor.UserID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return nil, domain.NewAccessDeniedError("manager %s cannot get daemons for application %s", actor.UserID, appID)
		}

		daemons, err := s.repo.GetDaemons(ctx, appID)
		if err != nil {
			return nil, err
		}

		return daemons, nil
	case admin.SystemRoleAdmin:
		daemons, err := s.repo.GetDaemons(ctx, appID)
		if err != nil {
			return nil, err
		}

		return daemons, nil
	case admin.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get daemons")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetDaemon implements the admin.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) GetDaemon(
	ctx context.Context,
	actor admin.Actor,
	appID, id admin.ID,
) (admin.Daemon, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Daemon{}, domain.NewAccessDeniedError("user %s cannot get daemon %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return admin.Daemon{}, domain.NewAccessDeniedError("manager %s cannot get daemon %s", actor.UserID, id)
		}

		daemon, err := s.repo.GetDaemon(ctx, appID, id)
		if err != nil {
			return admin.Daemon{}, err
		}

		return daemon, nil
	case admin.SystemRoleAdmin:
		user, err := s.repo.GetDaemon(ctx, appID, id)
		if err != nil {
			return admin.Daemon{}, err
		}

		return user, nil
	case admin.SystemRoleNone:
		return admin.Daemon{}, domain.NewAccessDeniedError("anonymous user cannot get daemon %s", id)
	default:
		return admin.Daemon{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateDaemon implements the admin.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) CreateDaemon(
	ctx context.Context,
	actor admin.Actor,
	daemon admin.Daemon,
) (admin.Daemon, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Daemon{}, domain.NewAccessDeniedError("user %s cannot create daemon", actor.UserID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != daemon.ApplicationID {
			return admin.Daemon{}, domain.NewAccessDeniedError("manager %s cannot create daemon", actor.UserID)
		}

		daemon.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateDaemon(ctx, daemon); err != nil {
			return admin.Daemon{}, err
		}

		return daemon, nil
	case admin.SystemRoleAdmin:
		daemon.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateDaemon(ctx, daemon); err != nil {
			return admin.Daemon{}, err
		}

		return daemon, nil
	case admin.SystemRoleNone:
		return admin.Daemon{}, domain.NewAccessDeniedError("anonymous user cannot create daemon")
	default:
		return admin.Daemon{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateDaemon implements the admin.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) UpdateDaemon(
	ctx context.Context,
	actor admin.Actor,
	daemon admin.Daemon,
) (admin.Daemon, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Daemon{}, domain.NewAccessDeniedError("user %s cannot update daemon %s", actor.UserID, daemon.ID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != daemon.ApplicationID {
			return admin.Daemon{}, domain.NewAccessDeniedError("manager %s cannot update daemon %s", actor.UserID, daemon.ID)
		}

		if err := s.repo.UpdateDaemon(ctx, daemon); err != nil {
			return admin.Daemon{}, err
		}

		return daemon, nil
	case admin.SystemRoleAdmin:
		if err := s.repo.UpdateDaemon(ctx, daemon); err != nil {
			return admin.Daemon{}, err
		}

		return daemon, nil
	case admin.SystemRoleNone:
		return admin.Daemon{}, domain.NewAccessDeniedError("anonymous user cannot update daemon %s", daemon.ID)
	default:
		return admin.Daemon{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteDaemon implements the admin.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) DeleteDaemon(
	ctx context.Context,
	actor admin.Actor,
	appID, id admin.ID,
) error {
	switch actor.Role {
	case admin.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete daemon %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.ApplicationID != appID {
			return domain.NewAccessDeniedError("manager %s cannot delete daemon %s", actor.UserID, id)
		}

		if err := s.repo.DeleteDaemon(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleAdmin:
		if err := s.repo.DeleteDaemon(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete daemon %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
