package service

import (
	"github.com/energimind/identity-server/internal/core/domain/auth"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
)

func toIdentityUser(ui oauth.UserInfo) auth.User {
	return auth.User{
		ID:         ui.ID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
