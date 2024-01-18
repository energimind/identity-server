package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewActor(t *testing.T) {
	t.Parallel()

	actor := NewActor("user1", "app1", SystemRoleAdmin)

	require.Equal(t, Actor{
		UserID:        "user1",
		ApplicationID: "app1",
		Role:          SystemRoleAdmin,
	}, actor)
}

func TestActor_IsValid(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		actor Actor
		valid bool
	}{
		"empty": {
			actor: Actor{},
			valid: false,
		},
		"invalid-role": {
			actor: Actor{
				UserID:        "user1",
				ApplicationID: "app1",
			},
			valid: false,
		},
		"invalid-userId": {
			actor: Actor{
				ApplicationID: "app1",
				Role:          SystemRoleAdmin,
			},
			valid: false,
		},
		"invalid-applicationId": {
			actor: Actor{
				UserID: "user1",
				Role:   SystemRoleAdmin,
			},
		},
		"valid": {
			actor: Actor{
				UserID:        "user1",
				ApplicationID: "app1",
				Role:          SystemRoleAdmin,
			},
			valid: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, test.valid, test.actor.IsValid())
		})
	}
}

func TestActor_String(t *testing.T) {
	t.Parallel()

	act := NewActor("user1", "app1", SystemRoleAdmin)

	require.Equal(t, "user1@app1[admin]", act.String())
}
