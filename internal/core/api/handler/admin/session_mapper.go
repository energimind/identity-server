package admin

import (
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/session"
)

func toSessionInfo(header session.Header, user admin.User) sessionInfo {
	return sessionInfo{
		ID:      header.SessionID,
		RealmID: header.RealmID,
		User:    toSessionUser(user),
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
