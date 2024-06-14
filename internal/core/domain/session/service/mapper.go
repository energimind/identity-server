package service

import (
	"github.com/energimind/identity-server/internal/core/domain/session"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
)

func toIdentityUser(ui oauth.UserInfo) session.User {
	return session.User{
		ID:         ui.ID,
		BindID:     ui.BindID,
		Name:       ui.Name,
		GivenName:  ui.GivenName,
		FamilyName: ui.FamilyName,
		Email:      ui.Email,
	}
}
