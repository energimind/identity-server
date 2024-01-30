package auth

import (
	"context"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
)

// ApplicationService is a service for managing applications.
//
// It implements the auth.ApplicationService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ApplicationService struct {
	repo  auth.ApplicationRepository
	idgen domain.IDGenerator
}

// NewApplicationService returns a new ApplicationService instance.
func NewApplicationService(
	repo auth.ApplicationRepository,
	idgen domain.IDGenerator,
) *ApplicationService {
	return &ApplicationService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the auth.ApplicationService interface.
var _ auth.ApplicationService = (*ApplicationService)(nil)

// GetApplications implements the auth.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) GetApplications(
	ctx context.Context,
	actor auth.Actor,
) ([]auth.Application, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get applications", actor.UserID)
	case auth.SystemRoleManager:
		app, err := s.repo.GetApplication(ctx, actor.ApplicationID)
		if err != nil {
			return nil, err
		}

		return []auth.Application{app}, nil
	case auth.SystemRoleAdmin:
		apps, err := s.repo.GetApplications(ctx)
		if err != nil {
			return nil, err
		}

		return apps, nil
	case auth.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get applications")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetApplication implements the auth.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) GetApplication(
	ctx context.Context,
	actor auth.Actor,
	id auth.ID,
) (auth.Application, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Application{}, domain.NewAccessDeniedError("user %s cannot get application %s", actor.UserID, id)
	case auth.SystemRoleManager:
		if actor.ApplicationID != id {
			return auth.Application{}, domain.NewAccessDeniedError("manager %s cannot get application %s", actor.UserID, id)
		}

		return s.repo.GetApplication(ctx, id)
	case auth.SystemRoleAdmin:
		return s.repo.GetApplication(ctx, id)
	case auth.SystemRoleNone:
		return auth.Application{}, domain.NewAccessDeniedError("anonymous user cannot get application %s", id)
	default:
		return auth.Application{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateApplication implements the auth.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) CreateApplication(
	ctx context.Context,
	actor auth.Actor,
	app auth.Application,
) (auth.Application, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Application{}, domain.NewAccessDeniedError("user %s cannot create application", actor.UserID)
	case auth.SystemRoleManager:
		return auth.Application{}, domain.NewAccessDeniedError("manager %s cannot create application", actor.UserID)
	case auth.SystemRoleAdmin:
		app.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateApplication(ctx, app); err != nil {
			return auth.Application{}, err
		}

		return app, nil
	case auth.SystemRoleNone:
		return auth.Application{}, domain.NewAccessDeniedError("anonymous user cannot create application")
	default:
		return auth.Application{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateApplication implements the auth.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) UpdateApplication(
	ctx context.Context,
	actor auth.Actor,
	app auth.Application,
) (auth.Application, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Application{}, domain.NewAccessDeniedError("user %s cannot update application %s", actor.UserID, app.ID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != app.ID {
			return auth.Application{}, domain.NewAccessDeniedError("manager %s cannot update application %s", actor.UserID, app.ID)
		}

		if err := s.repo.UpdateApplication(ctx, app); err != nil {
			return auth.Application{}, err
		}

		return app, nil
	case auth.SystemRoleAdmin:
		if err := s.repo.UpdateApplication(ctx, app); err != nil {
			return auth.Application{}, err
		}

		return app, nil
	case auth.SystemRoleNone:
		return auth.Application{}, domain.NewAccessDeniedError("anonymous user cannot update application %s", app.ID)
	default:
		return auth.Application{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteApplication implements the auth.ApplicationService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationService) DeleteApplication(
	ctx context.Context,
	actor auth.Actor,
	id auth.ID,
) error {
	switch actor.Role {
	case auth.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete application %s", actor.UserID, id)
	case auth.SystemRoleManager:
		return domain.NewAccessDeniedError("manager %s cannot delete application %s", actor.UserID, id)
	case auth.SystemRoleAdmin:
		return s.repo.DeleteApplication(ctx, id)
	case auth.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete application %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
