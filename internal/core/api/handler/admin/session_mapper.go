package admin

import (
	isclient "github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

func toSession(header isclient.Header, user admin.User) session {
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
