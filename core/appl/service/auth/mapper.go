package auth

import "github.com/energimind/identity-service/core/domain/auth"

func toUserInfo(ui userInfo) auth.UserInfo {
	return auth.UserInfo{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
