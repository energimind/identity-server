package login

import (
	"time"

	"golang.org/x/oauth2"
)

type session struct {
	ApplicationID string        `json:"applicationId"`
	Config        *config       `json:"config"`
	Token         *oauth2.Token `json:"token"`
	Timestamp     time.Time     `json:"timestamp"`
}

func newSession(applicationID string, config *config) *session {
	return &session{
		ApplicationID: applicationID,
		Config:        config,
		Timestamp:     time.Now(),
	}
}

func (s *session) updateToken(token *oauth2.Token) {
	s.Token = token
	s.Timestamp = time.Now()
}
