package service

import (
	"time"

	"github.com/energimind/identity-server/internal/core/domain/session"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
	"golang.org/x/oauth2"
)

type userSession struct {
	RealmID   string        `json:"realmId"`
	Config    *oauth.Config `json:"config"`
	Token     *oauth2.Token `json:"token,omitempty"`
	User      session.User  `json:"user,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

func newUserSession(realmID string, config *oauth.Config) *userSession {
	return &userSession{
		RealmID:   realmID,
		Config:    config,
		Timestamp: time.Now(),
	}
}

func (s *userSession) updateToken(token *oauth2.Token) {
	s.Token = token
	s.Timestamp = time.Now()
}

func (s *userSession) updateUser(user session.User) {
	s.User = user
	s.Timestamp = time.Now()
}
