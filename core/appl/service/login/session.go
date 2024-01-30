package login

import (
	"time"

	"golang.org/x/oauth2"
)

type session struct {
	Config    *config       `json:"config"`
	Token     *oauth2.Token `json:"token"`
	Timestamp time.Time     `json:"timestamp"`
}

func newSession(config *config) *session {
	return &session{
		Config:    config,
		Timestamp: time.Now(),
	}
}

func (s *session) updateToken(token *oauth2.Token) {
	s.Token = token
	s.Timestamp = time.Now()
}
