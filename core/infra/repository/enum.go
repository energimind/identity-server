package repository

const (
	dbProviderTypeNone dbProviderType = iota
	dbProviderTypeGoogle
)

const (
	dbSystemRoleNone dbSystemRole = iota
	dbSystemRoleUser
	dbSystemRoleManager
	dbSystemRoleAdmin
)

// All enums. Used for testing purposes to validate that all enum values are
// covered.
//
//nolint:gochecknoglobals,unused
var (
	allProviderTypes = []dbProviderType{dbProviderTypeNone, dbProviderTypeGoogle}
	allSystemRoles   = []dbSystemRole{dbSystemRoleNone, dbSystemRoleUser, dbSystemRoleManager, dbSystemRoleAdmin}
)

type dbProviderType int

type dbSystemRole int
