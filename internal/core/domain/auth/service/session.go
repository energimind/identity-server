package service

import (
	"time"

	"github.com/energimind/identity-server/internal/core/domain/auth"
	"github.com/energimind/identity-server/internal/core/infra/oauth"
	"golang.org/x/oauth2"
)

type userSession struct {
	ApplicationID string        `json:"applicationId"`
	Config        *oauth.Config `json:"config"`
	Token         *oauth2.Token `json:"token,omitempty"`
	User          auth.User     `json:"user,omitempty"`
	Timestamp     time.Time     `json:"timestamp"`
}

func newUserSession(applicationID string, config *oauth.Config) *userSession {
	return &userSession{
		ApplicationID: applicationID,
		Config:        config,
		Timestamp:     time.Now(),
	}
}

func (s *userSession) updateToken(token *oauth2.Token) {
	s.Token = token
	s.Timestamp = time.Now()
}

func (s *userSession) updateUser(user auth.User) {
	s.User = user
	s.Timestamp = time.Now()
}
