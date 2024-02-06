package handler

import "github.com/energimind/identity-service/core/domain/auth"

// adminActor is the actor for the admin role.
//
//nolint:gochecknoglobals // it is a constant
var adminActor = auth.Actor{Role: auth.SystemRoleAdmin}
