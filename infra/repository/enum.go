package repository

const dbProviderTypeGoogle dbProviderType = 1

const (
	dbSystemRoleUser dbSystemRole = iota + 1
	dbSystemRoleAdmin
	dbSystemRoleSysadmin
)

// All enums. Used for testing purposes to validate that all enum values are
// covered.
//
//nolint:gochecknoglobals
var (
	allProviderTypes = []dbProviderType{dbProviderTypeGoogle}
	allSystemRoles   = []dbSystemRole{dbSystemRoleUser, dbSystemRoleAdmin, dbSystemRoleSysadmin}
)

type dbProviderType int

type dbSystemRole int
