package auth

import (
	isclient "github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/domain/admin"
	"github.com/energimind/identity-server/internal/core/domain/auth"
)

func toClientSession(session auth.Session, user admin.User) isclient.Session {
	return isclient.Session{
		Header: isclient.Header{
			SessionID: session.Header.SessionID,
			RealmID:   session.Header.RealmID,
		},
		User: isclient.User{
			ID:          user.ID.String(),
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
		},
	}
}
