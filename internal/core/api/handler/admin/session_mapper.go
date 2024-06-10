package admin

import (
	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

func toInfo(session client.Session, user admin.User) info {
	return info{
		SessionInfo: toSessionInfo(session),
		UserInfo:    toUserInfo(user),
	}
}

func toSessionInfo(session client.Session) sessionInfo {
	return sessionInfo{
		SessionID:     session.SessionID,
		ApplicationID: session.ApplicationID,
	}
}

func toUserInfo(user admin.User) userInfo {
	return userInfo{
		ID:          user.ID.String(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role.String(),
	}
}
