package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// DaemonService is a service for managing daemons.
//
// It implements the service.DaemonService interface.
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

// Ensure service implements the service.DaemonService interface.
var _ admin.DaemonService = (*DaemonService)(nil)

// GetDaemons implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) GetDaemons(
	ctx context.Context,
	actor admin.Actor,
	realmID admin.ID,
) ([]admin.Daemon, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get daemons", actor.UserID)
	case admin.SystemRoleManager:
		if actor.RealmID != realmID {
			return nil, domain.NewAccessDeniedError("manager %s cannot get daemons for realm %s", actor.UserID, realmID)
		}

		daemons, err := s.repo.GetDaemons(ctx, realmID)
		if err != nil {
			return nil, err
		}

		return daemons, nil
	case admin.SystemRoleAdmin:
		daemons, err := s.repo.GetDaemons(ctx, realmID)
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

// GetDaemon implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) GetDaemon(
	ctx context.Context,
	actor admin.Actor,
	realmID, id admin.ID,
) (admin.Daemon, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Daemon{}, domain.NewAccessDeniedError("user %s cannot get daemon %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.RealmID != realmID {
			return admin.Daemon{}, domain.NewAccessDeniedError("manager %s cannot get daemon %s", actor.UserID, id)
		}

		daemon, err := s.repo.GetDaemon(ctx, realmID, id)
		if err != nil {
			return admin.Daemon{}, err
		}

		return daemon, nil
	case admin.SystemRoleAdmin:
		user, err := s.repo.GetDaemon(ctx, realmID, id)
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

// CreateDaemon implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) CreateDaemon(
	ctx context.Context,
	actor admin.Actor,
	daemon admin.Daemon,
) (admin.Daemon, error) {
	daemon, err := validateDaemon(daemon)
	if err != nil {
		return admin.Daemon{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Daemon{}, domain.NewAccessDeniedError("user %s cannot create daemon", actor.UserID)
	case admin.SystemRoleManager:
		if actor.RealmID != daemon.RealmID {
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

// UpdateDaemon implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) UpdateDaemon(
	ctx context.Context,
	actor admin.Actor,
	daemon admin.Daemon,
) (admin.Daemon, error) {
	daemon, err := validateDaemon(daemon)
	if err != nil {
		return admin.Daemon{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Daemon{}, domain.NewAccessDeniedError("user %s cannot update daemon %s", actor.UserID, daemon.ID)
	case admin.SystemRoleManager:
		if actor.RealmID != daemon.RealmID {
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

// DeleteDaemon implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) DeleteDaemon(
	ctx context.Context,
	actor admin.Actor,
	realmID, id admin.ID,
) error {
	switch actor.Role {
	case admin.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete daemon %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.RealmID != realmID {
			return domain.NewAccessDeniedError("manager %s cannot delete daemon %s", actor.UserID, id)
		}

		if err := s.repo.DeleteDaemon(ctx, realmID, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleAdmin:
		if err := s.repo.DeleteDaemon(ctx, realmID, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete daemon %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetAPIKeys implements the service.DaemonService interface.
func (s *DaemonService) GetAPIKeys(
	ctx context.Context,
	actor admin.Actor,
	realmID, daemonID admin.ID,
) ([]admin.APIKey, error) {
	daemon, err := s.GetDaemon(ctx, actor, realmID, daemonID)
	if err != nil {
		return nil, err
	}

	return daemon.APIKeys, nil
}

// GetAPIKey implements the service.DaemonService interface.
func (s *DaemonService) GetAPIKey(
	ctx context.Context,
	actor admin.Actor,
	realmID, daemonID, id admin.ID,
) (admin.APIKey, error) {
	daemon, err := s.GetDaemon(ctx, actor, realmID, daemonID)
	if err != nil {
		return admin.APIKey{}, err
	}

	for _, apiKey := range daemon.APIKeys {
		if apiKey.ID == id {
			return apiKey, nil
		}
	}

	return admin.APIKey{}, domain.NewNotFoundError("API key %s not found", id)
}

// CreateAPIKey implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) CreateAPIKey(
	ctx context.Context,
	actor admin.Actor,
	realmID, daemonID admin.ID,
	apiKey admin.APIKey,
) (admin.APIKey, error) {
	apiKey, err := validateAPIKey(apiKey)
	if err != nil {
		return admin.APIKey{}, err
	}

	daemon, err := s.GetDaemon(ctx, actor, realmID, daemonID)
	if err != nil {
		return admin.APIKey{}, err
	}

	apiKey.ID = admin.ID(s.idgen.GenerateID())

	daemon.APIKeys = append(daemon.APIKeys, apiKey)

	if uErr := s.repo.UpdateDaemon(ctx, daemon); uErr != nil {
		return admin.APIKey{}, uErr
	}

	return apiKey, nil
}

// UpdateAPIKey implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) UpdateAPIKey(
	ctx context.Context,
	actor admin.Actor,
	realmID, daemonID, id admin.ID,
	apiKey admin.APIKey,
) (admin.APIKey, error) {
	apiKey, err := validateAPIKey(apiKey)
	if err != nil {
		return admin.APIKey{}, err
	}

	daemon, err := s.GetDaemon(ctx, actor, realmID, daemonID)
	if err != nil {
		return admin.APIKey{}, err
	}

	for i, ak := range daemon.APIKeys {
		if ak.ID == id {
			daemon.APIKeys[i] = apiKey

			if uErr := s.repo.UpdateDaemon(ctx, daemon); uErr != nil {
				return admin.APIKey{}, uErr
			}

			return apiKey, nil
		}
	}

	return admin.APIKey{}, domain.NewNotFoundError("API key %s not found", id)
}

// DeleteAPIKey implements the service.DaemonService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *DaemonService) DeleteAPIKey(
	ctx context.Context,
	actor admin.Actor,
	realmID, daemonID, id admin.ID,
) error {
	daemon, err := s.GetDaemon(ctx, actor, realmID, daemonID)
	if err != nil {
		return err
	}

	var found bool

	for i, apiKey := range daemon.APIKeys {
		if apiKey.ID == id {
			daemon.APIKeys = append(daemon.APIKeys[:i], daemon.APIKeys[i+1:]...)
			found = true

			break
		}
	}

	if !found {
		return domain.NewNotFoundError("API key %s not found", id)
	}

	if uErr := s.repo.UpdateDaemon(ctx, daemon); uErr != nil {
		return uErr
	}

	return nil
}
