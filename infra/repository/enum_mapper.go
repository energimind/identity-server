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

func toKeyOwnerType(t auth.KeyOwnerType) dbKeyOwnerType {
	switch t {
	case auth.KeyOwnerTypeUser:
		return dbKeyOwnerTypeUser
	case auth.KeyOwnerTypeDaemon:
		return dbKeyOwnerTypeDaemon
	default:
		panic(fmt.Sprintf("unknown key owner type: %v", t))
	}
}

func fromKeyOwnerType(t dbKeyOwnerType) auth.KeyOwnerType {
	switch t {
	case dbKeyOwnerTypeUser:
		return auth.KeyOwnerTypeUser
	case dbKeyOwnerTypeDaemon:
		return auth.KeyOwnerTypeDaemon
	default:
		panic(fmt.Sprintf("unknown key owner type: %v", t))
	}
}
