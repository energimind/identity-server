package repository

import (
	"testing"

	"github.com/energimind/identity-service/domain/auth"
	"github.com/energimind/identity-service/test/utils"
	"github.com/stretchr/testify/require"
)

func Test_allEnumValuesAreMapped(t *testing.T) {
	t.Parallel()

	utils.CheckAllEnumValuesAreMapped(t, auth.AllProviderTypes, allProviderTypes, toProviderType)
	utils.CheckAllEnumValuesAreMapped(t, auth.AllSystemRoles, allSystemRoles, toSystemRole)

	utils.CheckAllEnumValuesAreMapped(t, allProviderTypes, auth.AllProviderTypes, fromProviderType)
	utils.CheckAllEnumValuesAreMapped(t, allSystemRoles, auth.AllSystemRoles, fromSystemRole)
}

func Test_enumMapperDefaultsOnInvalidEnum(t *testing.T) {
	require.Equal(t, dbProviderTypeNone, toProviderType("invalid"))
	require.Equal(t, dbSystemRoleNone, toSystemRole("invalid"))

	require.Equal(t, auth.ProviderTypeNone, fromProviderType(dbProviderType(-1)))
	require.Equal(t, auth.SystemRoleNone, fromSystemRole(dbSystemRole(-1)))
}
