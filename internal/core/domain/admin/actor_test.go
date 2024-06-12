package admin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewActor(t *testing.T) {
	t.Parallel()

	actor := NewActor("user1", "realm1", SystemRoleManager)

	require.Equal(t, Actor{
		UserID:  "user1",
		RealmID: "realm1",
		Role:    SystemRoleManager,
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
				UserID:  "user1",
				RealmID: "realm1",
			},
			valid: false,
		},
		"invalid-userId": {
			actor: Actor{
				RealmID: "realm1",
				Role:    SystemRoleManager,
			},
			valid: false,
		},
		"invalid-realmId": {
			actor: Actor{
				UserID: "user1",
				Role:   SystemRoleManager,
			},
		},
		"valid": {
			actor: Actor{
				UserID:  "user1",
				RealmID: "realm1",
				Role:    SystemRoleManager,
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

	act := NewActor("user1", "realm1", SystemRoleManager)

	require.Equal(t, "user1@realm1[manager]", act.String())
}
