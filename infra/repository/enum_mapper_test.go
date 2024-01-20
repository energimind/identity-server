package repository

import (
	"testing"

	"github.com/energimind/identity-service/domain/auth"
	"github.com/energimind/identity-service/test/utils"
	"github.com/stretchr/testify/require"
)

func Test_allEnumValuesAreMapped(t *testing.T) {
	utils.CheckAllEnumValuesAreMapped(t, auth.AllProviderTypes, allProviderTypes, toProviderType)
	utils.CheckAllEnumValuesAreMapped(t, auth.AllSystemRoles, allSystemRoles, toSystemRole)

	utils.CheckAllEnumValuesAreMapped(t, allProviderTypes, auth.AllProviderTypes, fromProviderType)
	utils.CheckAllEnumValuesAreMapped(t, allSystemRoles, auth.AllSystemRoles, fromSystemRole)
}

func Test_toProviderType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pt     auth.ProviderType
		exp    dbProviderType
		panics bool
	}{
		"google": {
			pt:  auth.ProviderTypeGoogle,
			exp: dbProviderTypeGoogle,
		},
		"unknown": {
			pt:     auth.ProviderType("unknown"),
			panics: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panics {
				require.Panics(t, func() {
					toProviderType(tc.pt)
				})

				return
			}

			require.Equal(t, tc.exp, toProviderType(tc.pt))
		})
	}
}

func Test_fromProviderType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pt     dbProviderType
		exp    auth.ProviderType
		panics bool
	}{
		"google": {
			pt:  dbProviderTypeGoogle,
			exp: auth.ProviderTypeGoogle,
		},
		"unknown": {
			pt:     dbProviderType(-1),
			panics: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panics {
				require.Panics(t, func() {
					fromProviderType(tc.pt)
				})

				return
			}

			require.Equal(t, tc.exp, fromProviderType(tc.pt))
		})
	}
}

func Test_toSystemRole(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		r      auth.SystemRole
		exp    dbSystemRole
		panics bool
	}{
		"user": {
			r:   auth.SystemRoleUser,
			exp: dbSystemRoleUser,
		},
		"admin": {
			r:   auth.SystemRoleManager,
			exp: dbSystemRoleManager,
		},
		"sysadmin": {
			r:   auth.SystemRoleAdmin,
			exp: dbSystemRoleAdmin,
		},
		"unknown": {
			r:      auth.SystemRole("unknown"),
			panics: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panics {
				require.Panics(t, func() {
					toSystemRole(tc.r)
				})

				return
			}

			require.Equal(t, tc.exp, toSystemRole(tc.r))
		})
	}
}

func Test_fromSystemRole(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		r      dbSystemRole
		exp    auth.SystemRole
		panics bool
	}{
		"user": {
			r:   dbSystemRoleUser,
			exp: auth.SystemRoleUser,
		},
		"admin": {
			r:   dbSystemRoleManager,
			exp: auth.SystemRoleManager,
		},
		"sysadmin": {
			r:   dbSystemRoleAdmin,
			exp: auth.SystemRoleAdmin,
		},
		"unknown": {
			r:      dbSystemRole(-1),
			panics: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.panics {
				require.Panics(t, func() {
					fromSystemRole(tc.r)
				})

				return
			}

			require.Equal(t, tc.exp, fromSystemRole(tc.r))
		})
	}
}
