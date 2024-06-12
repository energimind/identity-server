package admin

import (
	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

func toSession(si client.SessionInfo, user admin.User) session {
	return session{
		SessionID:     si.SessionID,
		ApplicationID: si.ApplicationID,
		User:          toSessionUser(user),
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
