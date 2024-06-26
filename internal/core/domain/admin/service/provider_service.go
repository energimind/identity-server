package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// ProviderService is a service for managing providers.
//
// It implements the service.ProviderService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type ProviderService struct {
	repo  admin.ProviderRepository
	idgen domain.IDGenerator
}

// NewProviderService returns a new ProviderService instance.
func NewProviderService(
	repo admin.ProviderRepository,
	idgen domain.IDGenerator,
) *ProviderService {
	return &ProviderService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the service.ProviderService interface.
var _ admin.ProviderService = (*ProviderService)(nil)

// GetProviders implements the service.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) GetProviders(
	ctx context.Context,
	actor admin.Actor,
) ([]admin.Provider, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get providers", actor.UserID)
	case admin.SystemRoleManager:
		return nil, domain.NewAccessDeniedError("manager %s cannot get providers", actor.UserID)
	case admin.SystemRoleAdmin:
		providers, err := s.repo.GetProviders(ctx)
		if err != nil {
			return nil, err
		}

		return providers, nil
	case admin.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get providers")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetProvider implements the service.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) GetProvider(
	ctx context.Context,
	actor admin.Actor,
	id admin.ID,
) (admin.Provider, error) {
	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Provider{}, domain.NewAccessDeniedError("user %s cannot get provider %s", actor.UserID, id)
	case admin.SystemRoleManager:
		return admin.Provider{}, domain.NewAccessDeniedError("manager %s cannot get provider %s", actor.UserID, id)
	case admin.SystemRoleAdmin:
		user, err := s.repo.GetProvider(ctx, id)
		if err != nil {
			return admin.Provider{}, err
		}

		return user, nil
	case admin.SystemRoleNone:
		return admin.Provider{}, domain.NewAccessDeniedError("anonymous user cannot get provider %s", id)
	default:
		return admin.Provider{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateProvider implements the service.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) CreateProvider(
	ctx context.Context,
	actor admin.Actor,
	provider admin.Provider,
) (admin.Provider, error) {
	provider, err := validateProvider(provider)
	if err != nil {
		return admin.Provider{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Provider{}, domain.NewAccessDeniedError("user %s cannot create provider", actor.UserID)
	case admin.SystemRoleManager:
		return admin.Provider{}, domain.NewAccessDeniedError("manager %s cannot create provider", actor.UserID)
	case admin.SystemRoleAdmin:
		provider.ID = admin.ID(s.idgen.GenerateID())

		if err := s.repo.CreateProvider(ctx, provider); err != nil {
			return admin.Provider{}, err
		}

		return provider, nil
	case admin.SystemRoleNone:
		return admin.Provider{}, domain.NewAccessDeniedError("anonymous user cannot create provider")
	default:
		return admin.Provider{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateProvider implements the service.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) UpdateProvider(
	ctx context.Context,
	actor admin.Actor,
	provider admin.Provider,
) (admin.Provider, error) {
	provider, err := validateProvider(provider)
	if err != nil {
		return admin.Provider{}, err
	}

	switch actor.Role {
	case admin.SystemRoleUser:
		return admin.Provider{}, domain.NewAccessDeniedError("user %s cannot update provider %s", actor.UserID, provider.ID)
	case admin.SystemRoleManager:
		return admin.Provider{}, domain.NewAccessDeniedError("manager %s cannot update provider %s", actor.UserID, provider.ID)
	case admin.SystemRoleAdmin:
		if err := s.repo.UpdateProvider(ctx, provider); err != nil {
			return admin.Provider{}, err
		}

		return provider, nil
	case admin.SystemRoleNone:
		return admin.Provider{}, domain.NewAccessDeniedError("anonymous user cannot update provider %s", provider.ID)
	default:
		return admin.Provider{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteProvider implements the service.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) DeleteProvider(
	ctx context.Context,
	actor admin.Actor,
	id admin.ID,
) error {
	switch actor.Role {
	case admin.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete provider %s", actor.UserID, id)
	case admin.SystemRoleManager:
		return domain.NewAccessDeniedError("manager %s cannot delete provider %s", actor.UserID, id)
	case admin.SystemRoleAdmin:
		if err := s.repo.DeleteProvider(ctx, id); err != nil {
			return err
		}

		return nil
	case admin.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete provider %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
