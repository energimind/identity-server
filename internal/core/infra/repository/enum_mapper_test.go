package repository

import (
	"testing"

	"github.com/energimind/go-kit/testutil/mapping"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/stretchr/testify/require"
)

func Test_allEnumValuesAreMapped(t *testing.T) {
	t.Parallel()

	mapping.CheckAllEnumValuesAreMapped(t, admin.AllProviderTypes, allProviderTypes, toProviderType)
	mapping.CheckAllEnumValuesAreMapped(t, admin.AllSystemRoles, allSystemRoles, toSystemRole)

	mapping.CheckAllEnumValuesAreMapped(t, allProviderTypes, admin.AllProviderTypes, fromProviderType)
	mapping.CheckAllEnumValuesAreMapped(t, allSystemRoles, admin.AllSystemRoles, fromSystemRole)
}

func Test_enumMapperDefaultsOnInvalidEnum(t *testing.T) {
	require.Equal(t, dbProviderTypeNone, toProviderType("invalid"))
	require.Equal(t, dbSystemRoleNone, toSystemRole("invalid"))

	require.Equal(t, admin.ProviderTypeNone, fromProviderType(dbProviderType(-1)))
	require.Equal(t, admin.SystemRoleNone, fromSystemRole(dbSystemRole(-1)))
}
