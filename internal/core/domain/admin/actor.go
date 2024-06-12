package admin

import (
	"fmt"
)

// Actor represents a user or realm that is performing an action.
// The actor is used in the context of a service call.
type Actor struct {
	UserID  ID
	RealmID ID
	Role    SystemRole
}

// NewActor returns a new actor.
func NewActor(userID ID, realmID ID, role SystemRole) Actor {
	return Actor{
		UserID:  userID,
		RealmID: realmID,
		Role:    role,
	}
}

// IsValid returns true if the actor is valid.
func (a Actor) IsValid() bool {
	return a.UserID != "" && a.RealmID != "" && a.Role != ""
}

// String implements the fmt.Stringer interface.
func (a Actor) String() string {
	return fmt.Sprintf("%s@%s[%s]", a.UserID, a.RealmID, a.Role)
}
