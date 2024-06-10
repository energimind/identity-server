package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// ApplicationLookupService provides a service to look up an application.
//
// It implements the service.ApplicationLookupService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ApplicationLookupService struct {
	applicationService admin.ApplicationService
	admin              admin.Actor
}

// NewApplicationLookupService returns a new ApplicationLookupService instance.
func NewApplicationLookupService(
	applicationService admin.ApplicationService,
) *ApplicationLookupService {
	return &ApplicationLookupService{
		applicationService: applicationService,
		admin:              admin.Actor{Role: admin.SystemRoleAdmin},
	}
}

// Ensure service implements the service.ApplicationLookupService interface.
var _ admin.ApplicationLookupService = (*ApplicationLookupService)(nil)

// LookupApplication implements the service.ApplicationLookupService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ApplicationLookupService) LookupApplication(
	ctx context.Context,
	appCode string,
) (admin.Application, error) {
	if appCode == "" {
		return admin.Application{}, domain.NewBadRequestError("application code must not be empty")
	}

	apps, err := s.applicationService.GetApplications(ctx, s.admin)
	if err != nil {
		return admin.Application{}, err
	}

	app, found := s.findProvider(apps, appCode)
	if !found {
		return admin.Application{}, domain.NewNotFoundError("application %s not found", appCode)
	}

	return app, nil
}

func (s *ApplicationLookupService) findProvider(apps []admin.Application, code string) (admin.Application, bool) {
	for _, app := range apps {
		if app.Code == code && app.Enabled {
			return app, true
		}
	}

	return admin.Application{}, false
}
