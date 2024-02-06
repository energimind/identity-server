package login

import "github.com/energimind/identity-service/core/domain/session"

func toUserInfo(ui userInfo) session.UserInfo {
	return session.UserInfo{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
