package repository

import (
	"fmt"

	"github.com/energimind/identity-service/domain/auth"
)

func toProviderType(t auth.ProviderType) dbProviderType {
	switch t {
	case auth.ProviderTypeGoogle:
		return dbProviderTypeGoogle
	default:
		panic(fmt.Sprintf("unknown provider type: %v", t))
	}
}

func fromProviderType(t dbProviderType) auth.ProviderType {
	switch t {
	case dbProviderTypeGoogle:
		return auth.ProviderTypeGoogle
	default:
		panic(fmt.Sprintf("unknown provider type: %v", t))
	}
}

func toSystemRole(r auth.SystemRole) dbSystemRole {
	switch r {
	case auth.SystemRoleUser:
		return dbSystemRoleUser
	case auth.SystemRoleAdmin:
		return dbSystemRoleAdmin
	case auth.SystemRoleSysadmin:
		return dbSystemRoleSysadmin
	default:
		panic(fmt.Sprintf("unknown system role: %v", r))
	}
}

func fromSystemRole(r dbSystemRole) auth.SystemRole {
	switch r {
	case dbSystemRoleUser:
		return auth.SystemRoleUser
	case dbSystemRoleAdmin:
		return auth.SystemRoleAdmin
	case dbSystemRoleSysadmin:
		return auth.SystemRoleSysadmin
	default:
		panic(fmt.Sprintf("unknown system role: %v", r))
	}
}
