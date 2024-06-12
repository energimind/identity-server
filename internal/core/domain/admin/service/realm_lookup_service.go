package service

import (
	"context"

	"github.com/energimind/identity-server/internal/core/domain"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

// RealmLookupService provides a service to look up a realm.
//
// It implements the service.RealmLookupService interface.
//
// We do not wrap the errors returned by the repository because they are already
// packed as domain errors. Therefore, we disable the wrapcheck linter for these calls.
type RealmLookupService struct {
	realmService admin.RealmService
	admin        admin.Actor
}

// NewRealmLookupService returns a new RealmLookupService instance.
func NewRealmLookupService(
	realmService admin.RealmService,
) *RealmLookupService {
	return &RealmLookupService{
		realmService: realmService,
		admin:        admin.Actor{Role: admin.SystemRoleAdmin},
	}
}

// Ensure service implements the service.RealmLookupService interface.
var _ admin.RealmLookupService = (*RealmLookupService)(nil)

// LookupRealm implements the service.RealmLookupService interface.
//
//nolint:wrapcheck // see comment in the header
func (s *RealmLookupService) LookupRealm(
	ctx context.Context,
	realmCode string,
) (admin.Realm, error) {
	if realmCode == "" {
		return admin.Realm{}, domain.NewBadRequestError("realm code must not be empty")
	}

	realms, err := s.realmService.GetRealms(ctx, s.admin)
	if err != nil {
		return admin.Realm{}, err
	}

	realm, found := s.findProvider(realms, realmCode)
	if !found {
		return admin.Realm{}, domain.NewNotFoundError("realm %s not found", realmCode)
	}

	return realm, nil
}

func (s *RealmLookupService) findProvider(realms []admin.Realm, code string) (admin.Realm, bool) {
	for _, realm := range realms {
		if realm.Code == code && realm.Enabled {
			return realm, true
		}
	}

	return admin.Realm{}, false
}
