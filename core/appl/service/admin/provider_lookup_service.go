package admin

import (
	"context"

	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
)

// ProviderLookupService provides a service to look up a provider
// from an application.
//
// It implements the admin.ProviderLookupService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ProviderLookupService struct {
	applicationService admin.ApplicationService
	providerService    admin.ProviderService
	admin              admin.Actor
}

// NewProviderLookupService returns a new ProviderLookupService instance.
func NewProviderLookupService(
	applicationService admin.ApplicationService,
	providerService admin.ProviderService,
) *ProviderLookupService {
	return &ProviderLookupService{
		applicationService: applicationService,
		providerService:    providerService,
		admin:              admin.Actor{Role: admin.SystemRoleAdmin},
	}
}

// Ensure service implements the admin.ProviderLookupService interface.
var _ admin.ProviderLookupService = (*ProviderLookupService)(nil)

// LookupProvider implements the admin.ProviderLookupService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderLookupService) LookupProvider(
	ctx context.Context,
	applicationCode, providerCode string,
) (admin.Provider, error) {
	if applicationCode == "" {
		return admin.Provider{}, domain.NewBadRequestError("application code must not be empty")
	}

	if providerCode == "" {
		return admin.Provider{}, domain.NewBadRequestError("provider code must not be empty")
	}

	apps, err := s.applicationService.GetApplications(ctx, s.admin)
	if err != nil {
		return admin.Provider{}, err
	}

	app, found := s.findApplication(apps, applicationCode)
	if !found {
		return admin.Provider{}, domain.NewNotFoundError("application %s not found", applicationCode)
	}

	providers, err := s.providerService.GetProviders(ctx, s.admin, app.ID)
	if err != nil {
		return admin.Provider{}, err
	}

	provider, found := s.findProvider(providers, providerCode)
	if !found {
		return admin.Provider{}, domain.NewNotFoundError("provider %s not found", providerCode)
	}

	return provider, nil
}

func (s *ProviderLookupService) findApplication(apps []admin.Application, code string) (admin.Application, bool) {
	for _, app := range apps {
		if app.Code == code && app.Enabled {
			return app, true
		}
	}

	return admin.Application{}, false
}

func (s *ProviderLookupService) findProvider(providers []admin.Provider, code string) (admin.Provider, bool) {
	for _, provider := range providers {
		if provider.Code == code && provider.Enabled {
			return provider, true
		}
	}

	return admin.Provider{}, false
}
