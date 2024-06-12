package service

import (
	isclient "github.com/energimind/identity-server/client"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
)

func toIdentityUser(ui oauth.UserInfo) isclient.User {
	return isclient.User{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
