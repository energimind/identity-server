package repository

import (
	"github.com/energimind/identity-service/domain/auth"
)

func toProviderType(t auth.ProviderType) dbProviderType {
	switch t {
	case auth.ProviderTypeNone:
		return dbProviderTypeNone
	case auth.ProviderTypeGoogle:
		return dbProviderTypeGoogle
	default:
		return dbProviderTypeNone
	}
}

func fromProviderType(t dbProviderType) auth.ProviderType {
	switch t {
	case dbProviderTypeNone:
		return auth.ProviderTypeNone
	case dbProviderTypeGoogle:
		return auth.ProviderTypeGoogle
	default:
		return auth.ProviderTypeNone
	}
}

func toSystemRole(r auth.SystemRole) dbSystemRole {
	switch r {
	case auth.SystemRoleNone:
		return dbSystemRoleNone
	case auth.SystemRoleUser:
		return dbSystemRoleUser
	case auth.SystemRoleManager:
		return dbSystemRoleManager
	case auth.SystemRoleAdmin:
		return dbSystemRoleAdmin
	default:
		return dbSystemRoleNone
	}
}

func fromSystemRole(r dbSystemRole) auth.SystemRole {
	switch r {
	case dbSystemRoleNone:
		return auth.SystemRoleNone
	case dbSystemRoleUser:
		return auth.SystemRoleUser
	case dbSystemRoleManager:
		return auth.SystemRoleManager
	case dbSystemRoleAdmin:
		return auth.SystemRoleAdmin
	default:
		return auth.SystemRoleNone
	}
}
