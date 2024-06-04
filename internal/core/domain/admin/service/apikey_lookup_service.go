package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// APIKeyLookupService provides a service for looking up API keys for
// a user or a daemon.
//
// It implements the service.APIKeyLookupService interface.
//
// We use the repository to look up the API key for a user and a daemon.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type APIKeyLookupService struct {
	userRepo   admin.UserRepository
	daemonRepo admin.DaemonRepository
}

// NewAPIKeyLookupService returns a new APIKeyLookupService instance.
func NewAPIKeyLookupService(
	userRepo admin.UserRepository,
	daemonRepo admin.DaemonRepository,
) *APIKeyLookupService {
	return &APIKeyLookupService{
		userRepo:   userRepo,
		daemonRepo: daemonRepo,
	}
}

// Ensure service implements the service.APIKeyLookupService interface.
var _ admin.APIKeyLookupService = (*APIKeyLookupService)(nil)

// LookupAPIKey implements the service.APIKeyLookupService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *APIKeyLookupService) LookupAPIKey(ctx context.Context, appID admin.ID, key string) (admin.APIKey, error) {
	fromUser, err := s.userRepo.GetAPIKey(ctx, appID, key)
	if err == nil {
		return fromUser, nil
	}

	fromDaemon, err := s.daemonRepo.GetAPIKey(ctx, appID, key)
	if err == nil {
		return fromDaemon, nil
	}

	return admin.APIKey{}, err
}
