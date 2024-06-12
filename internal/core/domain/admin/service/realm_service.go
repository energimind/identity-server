package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// RealmService is a service for managing realms.
//
// It implements the service.RealmService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type RealmService struct {
	repo  admin.RealmRepository
	idgen domain.IDGenerator
}

// NewRealmService returns a new RealmService instance.
func NewRealmService(
	repo admin.RealmRepository,
	idgen domain.IDGenerator,
) *RealmService {
	return &RealmService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the service.RealmService interface.
var _ admin.RealmService = (*RealmService)(nil)

// GetRealms implements the service.RealmService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *RealmService) GetRealms(
	ctx context.Context,
	actor admin.Actor,
) ([]admin.Realm, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get realms", actor.UserID)
	case admin.SystemRoleManager:
		realm, err := s.repo.GetRealm(ctx, actor.RealmID)
		if err != nil {
			return nil, err
		}

		return []admin.Realm{realm}, nil
	case admin.SystemRoleAdmin:
		realms, err := s.repo.GetRealms(ctx)
		if err != nil {
			return nil, err
		}

		return realms, nil
	case admin.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get realms")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetRealm implements the service.RealmService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *RealmService) GetRealm(
	ctx context.Context,
	actor admin.Actor,
	id admin.ID,
) (admin.Realm, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Realm{}, domain.NewAccessDeniedError("user %s cannot get realm %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.RealmID != id {
			return admin.Realm{}, domain.NewAccessDeniedError("manager %s cannot get realm %s", actor.UserID, id)
		}

		return s.repo.GetRealm(ctx, id)
	case admin.SystemRoleAdmin:
		return s.repo.GetRealm(ctx, id)
	case admin.SystemRoleNone:
		return admin.Realm{}, domain.NewAccessDeniedError("anonymous user cannot get realm %s", id)
	default:
		return admin.Realm{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateRealm implements the service.RealmService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *RealmService) CreateRealm(
	ctx context.Context,
	actor admin.Actor,
	realm admin.Realm,
) (admin.Realm, error) {
	realm, err := validateRealm(realm)
	if err != nil {
		return admin.Realm{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Realm{}, domain.NewAccessDeniedError("user %s cannot create realm", actor.UserID)
	case admin.SystemRoleManager:
		return admin.Realm{}, domain.NewAccessDeniedError("manager %s cannot create realm", actor.UserID)
	case admin.SystemRoleAdmin:
		realm.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateRealm(ctx, realm); err != nil {
			return admin.Realm{}, err
		}

		return realm, nil
	case admin.SystemRoleNone:
		return admin.Realm{}, domain.NewAccessDeniedError("anonymous user cannot create realm")
	default:
		return admin.Realm{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateRealm implements the service.RealmService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *RealmService) UpdateRealm(
	ctx context.Context,
	actor admin.Actor,
	realm admin.Realm,
) (admin.Realm, error) {
	realm, err := validateRealm(realm)
	if err != nil {
		return admin.Realm{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Realm{}, domain.NewAccessDeniedError("user %s cannot update realm %s", actor.UserID, realm.ID)
	case admin.SystemRoleManager:
		if actor.RealmID != realm.ID {
			return admin.Realm{}, domain.NewAccessDeniedError("manager %s cannot update realm %s", actor.UserID, realm.ID)
		}

		if err := s.repo.UpdateRealm(ctx, realm); err != nil {
			return admin.Realm{}, err
		}

		return realm, nil
	case admin.SystemRoleAdmin:
		if err := s.repo.UpdateRealm(ctx, realm); err != nil {
			return admin.Realm{}, err
		}

		return realm, nil
	case admin.SystemRoleNone:
		return admin.Realm{}, domain.NewAccessDeniedError("anonymous user cannot update realm %s", realm.ID)
	default:
		return admin.Realm{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteRealm implements the service.RealmService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *RealmService) DeleteRealm(
	ctx context.Context,
	actor admin.Actor,
	id admin.ID,
) error {
	switch actor.Role {
	case admin.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete realm %s", actor.UserID, id)
	case admin.SystemRoleManager:
		return domain.NewAccessDeniedError("manager %s cannot delete realm %s", actor.UserID, id)
	case admin.SystemRoleAdmin:
		return s.repo.DeleteRealm(ctx, id)
	case admin.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete realm %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
