package admin

import (
	"strings"
)

const fieldSeparator = ":"

// UserSession is a struct that contains user session information.
type UserSession struct {
	SessionID     string
	ApplicationID string
	UserID        string
	UserRole      string
}

// NewUserSession creates a new UserSession with the given parameters.
func NewUserSession(sessionID, applicationID, userID, userRole string) UserSession {
	return UserSession{
		SessionID:     sessionID,
		ApplicationID: applicationID,
		UserID:        userID,
		UserRole:      userRole,
	}
}

// Serialize returns the serialized representation of the UserSession.
// This is useful for storing the session in a cookie.
func (s UserSession) Serialize() string {
	return s.SessionID + fieldSeparator +
		s.ApplicationID + fieldSeparator +
		s.UserID + fieldSeparator +
		s.UserRole
}

// DeserializeUserSession deserializes the given string into a UserSession.
func DeserializeUserSession(serialized string) (UserSession, error) {
	const expectedPartCount = 4

	parts := strings.Split(serialized, fieldSeparator)

	if len(parts) != expectedPartCount {
		return UserSession{}, NewUserSessionError("invalid serialized user session")
	}

	return UserSession{
		SessionID:     parts[0],
		ApplicationID: parts[1],
		UserID:        parts[2],
		UserRole:      parts[3],
	}, nil
}
