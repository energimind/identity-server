package admin

import (
	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/domain/admin"
)

func toInfo(sessionInfo client.SessionInfo, user admin.User) info {
	return info{
		SessionInfo: toSessionInfo(sessionInfo),
		UserInfo:    toUserInfo(user),
	}
}

func toSessionInfo(session client.SessionInfo) sessionInfo {
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
