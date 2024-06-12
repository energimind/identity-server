package admin

import (
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/auth"
)

func toSession(header auth.Header, user admin.User) session {
	return session{
		ID:    header.SessionID,
		AppID: header.ApplicationID,
		User:  toSessionUser(user),
	}
}

func toSessionUser(au admin.User) sessionUser {
	return sessionUser{
		ID:          au.ID.String(),
		Username:    au.Username,
		DisplayName: au.DisplayName,
		Email:       au.Email,
		Role:        au.Role.String(),
	}
}
