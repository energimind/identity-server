package login

import "github.com/energimind/identity-service/core/domain/login"

func toUserInfo(ui userInfo) login.UserInfo {
	return login.UserInfo{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
