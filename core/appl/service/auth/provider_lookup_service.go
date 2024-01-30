package auth

import (
	"context"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/auth"
)

// ProviderLookupService provides a service to lookup a provider
// from an application.
//
// It implements the auth.ProviderLookupService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ProviderLookupService struct {
	applicationService auth.ApplicationService
	providerService    auth.ProviderService
	admin              auth.Actor
}

// NewProviderLookupService returns a new ProviderLookupService instance.
func NewProviderLookupService(
	applicationService auth.ApplicationService,
	providerService auth.ProviderService,
) *ProviderLookupService {
	return &ProviderLookupService{
		applicationService: applicationService,
		providerService:    providerService,
		admin:              auth.Actor{Role: auth.SystemRoleAdmin},
	}
}

// Ensure service implements the auth.ProviderLookupService interface.
var _ auth.ProviderLookupService = (*ProviderLookupService)(nil)

// LookupProvider implements the auth.ProviderLookupService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderLookupService) LookupProvider(
	ctx context.Context,
	applicationCode, providerCode string,
) (auth.Provider, error) {
	if applicationCode == "" {
		return auth.Provider{}, domain.NewBadRequestError("application code must not be empty")
	}

	if providerCode == "" {
		return auth.Provider{}, domain.NewBadRequestError("provider code must not be empty")
	}

	apps, err := s.applicationService.GetApplications(ctx, s.admin)
	if err != nil {
		return auth.Provider{}, err
	}

	app, found := s.findApplication(apps, applicationCode)
	if !found {
		return auth.Provider{}, domain.NewNotFoundError("application %s not found", applicationCode)
	}

	providers, err := s.providerService.GetProviders(ctx, s.admin, app.ID)
	if err != nil {
		return auth.Provider{}, err
	}

	provider, found := s.findProvider(providers, providerCode)
	if !found {
		return auth.Provider{}, domain.NewNotFoundError("provider %s not found", providerCode)
	}

	return provider, nil
}

func (s *ProviderLookupService) findApplication(apps []auth.Application, code string) (auth.Application, bool) {
	for _, app := range apps {
		if app.Code == code && app.Enabled {
			return app, true
		}
	}

	return auth.Application{}, false
}

func (s *ProviderLookupService) findProvider(providers []auth.Provider, code string) (auth.Provider, bool) {
	for _, provider := range providers {
		if provider.Code == code && provider.Enabled {
			return provider, true
		}
	}

	return auth.Provider{}, false
}
