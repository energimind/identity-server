package auth

import (
	"github.com/energimind/identity-server/core/domain/auth"
	"github.com/energimind/identity-server/core/infra/oauth"
)

func toUserInfo(ui oauth.UserInfo) auth.UserInfo {
	return auth.UserInfo{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
