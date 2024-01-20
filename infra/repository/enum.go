package repository

const dbProviderTypeGoogle dbProviderType = 1

const (
	dbSystemRoleUser dbSystemRole = iota + 1
	dbSystemRoleManager
	dbSystemRoleAdmin
)

// All enums. Used for testing purposes to validate that all enum values are
// covered.
//
//nolint:gochecknoglobals,unused
var (
	allProviderTypes = []dbProviderType{dbProviderTypeGoogle}
	allSystemRoles   = []dbSystemRole{dbSystemRoleUser, dbSystemRoleManager, dbSystemRoleAdmin}
)

type dbProviderType int

type dbSystemRole int
