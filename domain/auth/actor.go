package auth

import (
	"fmt"
)

// Actor represents a user or application that is performing an action.
// The actor is used in the context of a service call.
type Actor struct {
	UserID        string
	ApplicationID string
	Role          SystemRole
}

// NewActor returns a new actor.
func NewActor(userID string, applicationID string, role SystemRole) Actor {
	return Actor{
		UserID:        userID,
		ApplicationID: applicationID,
		Role:          role,
	}
}

// IsValid returns true if the actor is valid.
func (a Actor) IsValid() bool {
	return a.UserID != "" && a.ApplicationID != "" && a.Role != ""
}

// String implements the fmt.Stringer interface.
func (a Actor) String() string {
	return fmt.Sprintf("%s@%s[%s]", a.UserID, a.ApplicationID, a.Role)
}
