package service

import (
	"github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
)

func toUserInfo(ui oauth.UserInfo) client.UserInfo {
	return client.UserInfo{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
