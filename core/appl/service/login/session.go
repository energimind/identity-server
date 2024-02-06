package login

import (
	"time"

	"golang.org/x/oauth2"
)

type userSession struct {
	ApplicationID string        `json:"applicationId"`
	Config        *config       `json:"config"`
	Token         *oauth2.Token `json:"token"`
	Timestamp     time.Time     `json:"timestamp"`
}

func newUserSession(applicationID string, config *config) *userSession {
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
