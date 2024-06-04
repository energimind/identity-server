package repository

import (
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

func toProviderType(t admin.ProviderType) dbProviderType {
	switch t {
	case admin.ProviderTypeNone:
		return dbProviderTypeNone
	case admin.ProviderTypeGoogle:
		return dbProviderTypeGoogle
	default:
		return dbProviderTypeNone
	}
}

func fromProviderType(t dbProviderType) admin.ProviderType {
	switch t {
	case dbProviderTypeNone:
		return admin.ProviderTypeNone
	case dbProviderTypeGoogle:
		return admin.ProviderTypeGoogle
	default:
		return admin.ProviderTypeNone
	}
}

func toSystemRole(r admin.SystemRole) dbSystemRole {
	switch r {
	case admin.SystemRoleNone:
		return dbSystemRoleNone
	case admin.SystemRoleUser:
		return dbSystemRoleUser
	case admin.SystemRoleManager:
		return dbSystemRoleManager
	case admin.SystemRoleAdmin:
		return dbSystemRoleAdmin
	default:
		return dbSystemRoleNone
	}
}

func fromSystemRole(r dbSystemRole) admin.SystemRole {
	switch r {
	case dbSystemRoleNone:
		return admin.SystemRoleNone
	case dbSystemRoleUser:
		return admin.SystemRoleUser
	case dbSystemRoleManager:
		return admin.SystemRoleManager
	case dbSystemRoleAdmin:
		return admin.SystemRoleAdmin
	default:
		return admin.SystemRoleNone
	}
}
