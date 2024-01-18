package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPrincipal(t *testing.T) {
	t.Parallel()

	principal := NewPrincipal("user1", "app1")

	require.Equal(t, Principal{
		UserID:        "user1",
		ApplicationID: "app1",
	}, principal)
}

func TestPrincipal_ForRole(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		principal Principal
		role      SystemRole
		expected  Principal
	}{
		"user": {
			principal: Principal{
				UserID:        "user1",
				ApplicationID: "app1",
			},
			role: SystemRoleUser,
			expected: Principal{
				UserID:        "user1",
				ApplicationID: "app1",
			},
		},
		"admin": {
			principal: Principal{
				UserID:        "user1",
				ApplicationID: "app1",
			},
			role: SystemRoleAdmin,
			expected: Principal{
				ApplicationID: "app1",
			},
		},
		"sysadmin": {
			principal: Principal{
				UserID:        "user1",
				ApplicationID: "app1",
			},
			role:     SystemRoleSysadmin,
			expected: Principal{},
		},
		"invalidRole": {
			principal: Principal{
				UserID:        "user1",
				ApplicationID: "app1",
			},
			role:     SystemRole("invalid"),
			expected: prohibitedPrincipal(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			principal := test.principal.ForRole(test.role)

			require.Equal(t, test.expected, principal)
		})
	}
}

func TestPrincipal_String(t *testing.T) {
	t.Parallel()

	t.Run("all", func(t *testing.T) {
		principal := Principal{
			UserID:        "user1",
			ApplicationID: "app1",
		}

		require.Equal(t, "user1@app1", principal.String())
	})

	t.Run("empty", func(t *testing.T) {
		principal := Principal{}

		require.Equal(t, "*@*", principal.String())
	})
}

func Test_prohibitedPrincipal(t *testing.T) {
	t.Parallel()

	require.Equal(t, Principal{
		UserID:        "0",
		ApplicationID: "0",
	}, prohibitedPrincipal())
}
