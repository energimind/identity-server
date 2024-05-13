// Package local supports the local admin login.
// This admin login is used for bootstrapping the system.
package local

import "github.com/energimind/identity-server/core/domain/admin"

// Local admin constants.
const (
	AdminProviderCode  = "local"
	AdminApplicationID = "localApplicationId"
	AdminSessionID     = "localSessionId"
	AdminID            = "localAdminId"
	AdminRole          = admin.SystemRoleAdmin
)
