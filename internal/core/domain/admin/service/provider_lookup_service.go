package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// ProviderLookupService provides a service to look up a provider
// from an application.
//
// It implements the service.ProviderLookupService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ProviderLookupService struct {
	providerService admin.ProviderService
	admin           admin.Actor
}

// NewProviderLookupService returns a new ProviderLookupService instance.
func NewProviderLookupService(
	providerService admin.ProviderService,
) *ProviderLookupService {
	return &ProviderLookupService{
		providerService: providerService,
		admin:           admin.Actor{Role: admin.SystemRoleAdmin},
	}
}

// Ensure service implements the service.ProviderLookupService interface.
var _ admin.ProviderLookupService = (*ProviderLookupService)(nil)

// LookupProvider implements the service.ProviderLookupService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderLookupService) LookupProvider(
	ctx context.Context,
	providerCode string,
) (admin.Provider, error) {
	if providerCode == "" {
		return admin.Provider{}, domain.NewBadRequestError("provider code must not be empty")
	}

	providers, err := s.providerService.GetProviders(ctx, s.admin)
	if err != nil {
		return admin.Provider{}, err
	}

	provider, found := s.findProvider(providers, providerCode)
	if !found {
		return admin.Provider{}, domain.NewNotFoundError("provider %s not found", providerCode)
	}

	return provider, nil
}

func (s *ProviderLookupService) findProvider(providers []admin.Provider, code string) (admin.Provider, bool) {
	for _, provider := range providers {
		if provider.Code == code && provider.Enabled {
			return provider, true
		}
	}

	return admin.Provider{}, false
}
