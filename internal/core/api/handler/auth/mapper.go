package auth

import "github.com/energimind/identity-server/internal/core/domain/auth"

// toInfo converts an auth.Info to an info.
func toInfo(ai auth.Info) info {
	return info{
		SessionInfo: sessionInfo{
			SessionID:     ai.SessionInfo.SessionID,
			ApplicationID: ai.SessionInfo.ApplicationID,
		},
		UserInfo: userInfo{
			ID:         ai.UserInfo.ID,
			Name:       ai.UserInfo.Name,
			GivenName:  ai.UserInfo.GivenName,
			FamilyName: ai.UserInfo.FamilyName,
			Email:      ai.UserInfo.Email,
		},
	}
}
