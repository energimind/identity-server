package service

import (
	"context"

	"github.com/energimind/identity-service/domain"
	"github.com/energimind/identity-service/domain/auth"
)

// ProviderService is a service for managing providers.
//
// It implements the auth.ProviderService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
//
// Some methods are reported as to complex by the linter. We disable the linter for
// these methods, because they are not too complex, but just have a lot of error handling.
type ProviderService struct {
	repo  auth.ProviderRepository
	idgen domain.IDGenerator
}

// NewProviderService returns a new ProviderService instance.
func NewProviderService(
	repo auth.ProviderRepository,
	idgen domain.IDGenerator,
) *ProviderService {
	return &ProviderService{
		repo:  repo,
		idgen: idgen,
	}
}

// Ensure service implements the auth.ProviderService interface.
var _ auth.ProviderService = (*ProviderService)(nil)

// GetProviders implements the auth.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) GetProviders(
	ctx context.Context,
	actor auth.Actor,
	appID auth.ID,
) ([]auth.Provider, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return nil, domain.NewAccessDeniedError("user %s cannot get providers", actor.UserID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != appID {
			return nil, domain.NewAccessDeniedError("manager %s cannot get providers for application %s", actor.UserID, appID)
		}

		providers, err := s.repo.GetProviders(ctx, appID)
		if err != nil {
			return nil, err
		}

		return providers, nil
	case auth.SystemRoleAdmin:
		providers, err := s.repo.GetProviders(ctx, appID)
		if err != nil {
			return nil, err
		}

		return providers, nil
	case auth.SystemRoleNone:
		return nil, domain.NewAccessDeniedError("anonymous user cannot get providers")
	default:
		return nil, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// GetProvider implements the auth.ProviderService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *ProviderService) GetProvider(
	ctx context.Context,
	actor auth.Actor,
	appID, id auth.ID,
) (auth.Provider, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Provider{}, domain.NewAccessDeniedError("user %s cannot get provider %s", actor.UserID, id)
	case auth.SystemRoleManager:
		provider, err := s.repo.GetProvider(ctx, id)
		if err != nil {
			return auth.Provider{}, err
		}

		if actor.ApplicationID != appID {
			return auth.Provider{}, domain.NewAccessDeniedError("manager %s cannot get provider %s", actor.UserID, id)
		}

		return provider, nil
	case auth.SystemRoleAdmin:
		user, err := s.repo.GetProvider(ctx, id)
		if err != nil {
			return auth.Provider{}, err
		}

		return user, nil
	case auth.SystemRoleNone:
		return auth.Provider{}, domain.NewAccessDeniedError("anonymous user cannot get provider %s", id)
	default:
		return auth.Provider{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// CreateProvider implements the auth.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) CreateProvider(
	ctx context.Context,
	actor auth.Actor,
	provider auth.Provider,
) (auth.Provider, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Provider{}, domain.NewAccessDeniedError("user %s cannot create provider", actor.UserID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != provider.ApplicationID {
			return auth.Provider{}, domain.NewAccessDeniedError("manager %s cannot create provider", actor.UserID)
		}

		provider.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateProvider(ctx, provider); err != nil {
			return auth.Provider{}, err
		}

		return provider, nil
	case auth.SystemRoleAdmin:
		provider.ID = auth.ID(s.idgen.GenerateID())

		if err := s.repo.CreateProvider(ctx, provider); err != nil {
			return auth.Provider{}, err
		}

		return provider, nil
	case auth.SystemRoleNone:
		return auth.Provider{}, domain.NewAccessDeniedError("anonymous user cannot create provider")
	default:
		return auth.Provider{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// UpdateProvider implements the auth.ProviderService interface.
//
//nolint:wrapcheck,cyclop // see comment in the header
func (s *ProviderService) UpdateProvider(
	ctx context.Context,
	actor auth.Actor,
	provider auth.Provider,
) (auth.Provider, error) {
	switch actor.Role {
	case auth.SystemRoleUser:
		return auth.Provider{}, domain.NewAccessDeniedError("user %s cannot update provider %s", actor.UserID, provider.ID)
	case auth.SystemRoleManager:
		if actor.ApplicationID != provider.ApplicationID {
			return auth.Provider{}, domain.NewAccessDeniedError("manager %s cannot update provider %s", actor.UserID, provider.ID)
		}

		if err := s.repo.UpdateProvider(ctx, provider); err != nil {
			return auth.Provider{}, err
		}

		return provider, nil
	case auth.SystemRoleAdmin:
		if err := s.repo.UpdateProvider(ctx, provider); err != nil {
			return auth.Provider{}, err
		}

		return provider, nil
	case auth.SystemRoleNone:
		return auth.Provider{}, domain.NewAccessDeniedError("anonymous user cannot update provider %s", provider.ID)
	default:
		return auth.Provider{}, domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}

// DeleteProvider implements the auth.ProviderService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *ProviderService) DeleteProvider(
	ctx context.Context,
	actor auth.Actor,
	appID, id auth.ID,
) error {
	switch actor.Role {
	case auth.SystemRoleUser:
		return domain.NewAccessDeniedError("user %s cannot delete provider %s", actor.UserID, id)
	case auth.SystemRoleManager:
		if actor.ApplicationID != appID {
			return domain.NewAccessDeniedError("manager %s cannot delete provider %s", actor.UserID, id)
		}

		if err := s.repo.DeleteProvider(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case auth.SystemRoleAdmin:
		if err := s.repo.DeleteProvider(ctx, appID, id); err != nil {
			return err
		}

		return nil
	case auth.SystemRoleNone:
		return domain.NewAccessDeniedError("anonymous user cannot delete provider %s", id)
	default:
		return domain.NewAccessDeniedError("unknown actor role %s", actor.Role)
	}
}
