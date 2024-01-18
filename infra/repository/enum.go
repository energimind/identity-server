package repository

const dbProviderTypeGoogle dbProviderType = 1

const (
	dbSystemRoleUser dbSystemRole = iota + 1
	dbSystemRoleAdmin
	dbSystemRoleSysadmin
)

const (
	dbKeyOwnerTypeUser dbKeyOwnerType = iota + 1
	dbKeyOwnerTypeDaemon
)

// All enums. Used for testing purposes to validate that all enum values are
// covered.
//
//nolint:gochecknoglobals
var (
	allProviderTypes = []dbProviderType{dbProviderTypeGoogle}
	allSystemRoles   = []dbSystemRole{dbSystemRoleUser, dbSystemRoleAdmin, dbSystemRoleSysadmin}
	allKeyOwnerTypes = []dbKeyOwnerType{dbKeyOwnerTypeUser, dbKeyOwnerTypeDaemon}
)

type dbProviderType int

type dbSystemRole int

type dbKeyOwnerType int
