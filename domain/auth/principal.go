package auth

import (
	"fmt"
)

// Principal represents the principal that is performing the repository action.
// Its fields are used to narrow the scope of the repository action.
// Any or all of the fields can be empty, indicating that the action is not
// scoped to that field.
// The principal is deduced from the user's system role.
type Principal struct {
	UserID        ID
	ApplicationID ID
}

// NewPrincipal returns a new principal.
func NewPrincipal(userID ID, applicationID ID) Principal {
	return Principal{
		UserID:        userID,
		ApplicationID: applicationID,
	}
}

// ForRole returns a new principal for the given system role.
func (p Principal) ForRole(r SystemRole) Principal {
	switch r {
	case SystemRoleUser: // keep narrow
		return p
	case SystemRoleAdmin: // scope to application
		return Principal{
			ApplicationID: p.ApplicationID,
		}
	case SystemRoleSysadmin: // scope to nothing
		return Principal{}
	default:
		return prohibitedPrincipal()
	}
}

// String implements the fmt.Stringer interface.
func (p Principal) String() string {
	star := func(s string) string {
		if s == "" {
			return "*"
		}

		return s
	}

	return fmt.Sprintf("%s@%s", star(string(p.UserID)), star(string(p.ApplicationID)))
}

func prohibitedPrincipal() Principal {
	const noAccess = "0"

	return Principal{
		UserID:        noAccess,
		ApplicationID: noAccess,
	}
}
