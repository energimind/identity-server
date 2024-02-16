package admin

import (
	"context"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
)

// ApplicationService is a service for managing applications.
//
// It implements the admin.ApplicationService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ApplicationService struct {
	repo  admin.ApplicationRepository
	idgen domain.IDGenerator
}

// NewApplicationService returns a new ApplicationService instance.
func NewApplicationService(
	repo admin.ApplicationRepository,
	idgen domain.IDGenerator,
) *ApplicationService {
	return &ApplicationService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the admin.ApplicationService interface.
var _ admin.ApplicationService = (*ApplicationService)(nil)

// GetApplications implements the admin.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) GetApplications(
	ctx context.Context,
	actor admin.Actor,
) ([]admin.Application, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get applications", actor.UserID)
	case admin.SystemRoleManager:
		app, err := s.repo.GetApplication(ctx, actor.ApplicationID)
		if err != nil {
			return nil, err
		}

		return []admin.Application{app}, nil
	case admin.SystemRoleAdmin:
		apps, err := s.repo.GetApplications(ctx)
		if err != nil {
			return nil, err
		}

		return apps, nil
	case admin.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get applications")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetApplication implements the admin.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) GetApplication(
	ctx context.Context,
	actor admin.Actor,
	id admin.ID,
) (admin.Application, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Application{}, domain.NewAccessDeniedError("user %s cannot get application %s", actor.UserID, id)
	case admin.SystemRoleManager:
		if actor.ApplicationID != id {
			return admin.Application{}, domain.NewAccessDeniedError("manager %s cannot get application %s", actor.UserID, id)
		}

		return s.repo.GetApplication(ctx, id)
	case admin.SystemRoleAdmin:
		return s.repo.GetApplication(ctx, id)
	case admin.SystemRoleNone:
		return admin.Application{}, domain.NewAccessDeniedError("anonymous user cannot get application %s", id)
	default:
		return admin.Application{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateApplication implements the admin.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) CreateApplication(
	ctx context.Context,
	actor admin.Actor,
	app admin.Application,
) (admin.Application, error) {
	app, err := validateApplication(app)
	if err != nil {
		return admin.Application{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Application{}, domain.NewAccessDeniedError("user %s cannot create application", actor.UserID)
	case admin.SystemRoleManager:
		return admin.Application{}, domain.NewAccessDeniedError("manager %s cannot create application", actor.UserID)
	case admin.SystemRoleAdmin:
		app.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateApplication(ctx, app); err != nil {
			return admin.Application{}, err
		}

		return app, nil
	case admin.SystemRoleNone:
		return admin.Application{}, domain.NewAccessDeniedError("anonymous user cannot create application")
	default:
		return admin.Application{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateApplication implements the admin.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) UpdateApplication(
	ctx context.Context,
	actor admin.Actor,
	app admin.Application,
) (admin.Application, error) {
	app, err := validateApplication(app)
	if err != nil {
		return admin.Application{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Application{}, domain.NewAccessDeniedError("user %s cannot update application %s", actor.UserID, app.ID)
	case admin.SystemRoleManager:
		if actor.ApplicationID != app.ID {
			return admin.Application{}, domain.NewAccessDeniedError("manager %s cannot update application %s", actor.UserID, app.ID)
		}

		if err := s.repo.UpdateApplication(ctx, app); err != nil {
			return admin.Application{}, err
		}

		return app, nil
	case admin.SystemRoleAdmin:
		if err := s.repo.UpdateApplication(ctx, app); err != nil {
			return admin.Application{}, err
		}

		return app, nil
	case admin.SystemRoleNone:
		return admin.Application{}, domain.NewAccessDeniedError("anonymous user cannot update application %s", app.ID)
	default:
		return admin.Application{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteApplication implements the admin.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) DeleteApplication(
	ctx context.Context,
	actor admin.Actor,
	id admin.ID,
) error {
	switch actor.Role {
	case admin.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete application %s", actor.UserID, id)
	case admin.SystemRoleManager:
		return domain.NewAccessDeniedError("manager %s cannot delete application %s", actor.UserID, id)
	case admin.SystemRoleAdmin:
		return s.repo.DeleteApplication(ctx, id)
	case admin.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete application %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
