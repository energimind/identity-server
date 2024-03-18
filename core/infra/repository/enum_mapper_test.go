package repository

import (
	"testing"

	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/energimind/identity-server/core/testutil"
	"github.com/stretchr/testify/require"
)

func Test_allEnumValuesAreMapped(t *testing.T) {
	t.Parallel()

	testutil.CheckAllEnumValuesAreMapped(t, admin.AllProviderTypes, allProviderTypes, toProviderType)
	testutil.CheckAllEnumValuesAreMapped(t, admin.AllSystemRoles, allSystemRoles, toSystemRole)

	testutil.CheckAllEnumValuesAreMapped(t, allProviderTypes, admin.AllProviderTypes, fromProviderType)
	testutil.CheckAllEnumValuesAreMapped(t, allSystemRoles, admin.AllSystemRoles, fromSystemRole)
}

func Test_enumMapperDefaultsOnInvalidEnum(t *testing.T) {
	require.Equal(t, dbProviderTypeNone, toProviderType("invalid"))
	require.Equal(t, dbSystemRoleNone, toSystemRole("invalid"))

	require.Equal(t, admin.ProviderTypeNone, fromProviderType(dbProviderType(-1)))
	require.Equal(t, admin.SystemRoleNone, fromSystemRole(dbSystemRole(-1)))
}
