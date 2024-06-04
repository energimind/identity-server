package sessioncookie

import (
	"strings"

	"github.com/energimind/identity-server/internal/core/domain"
)

const fieldSeparator = ":"

// serializeUserSession returns the serialized representation of the UserSession.
func serializeUserSession(us domain.UserSession) string {
	return us.SessionID + fieldSeparator +
		us.ApplicationID + fieldSeparator +
		us.UserID + fieldSeparator +
		us.UserRole
}

// deserializeUserSession deserializes the given string into a UserSession.
func deserializeUserSession(serialized string) (domain.UserSession, error) {
	const expectedPartCount = 4

	parts := strings.Split(serialized, fieldSeparator)

	if len(parts) != expectedPartCount {
		return domain.UserSession{}, NewError("invalid serialized user session")
	}

	return domain.UserSession{
		SessionID:     parts[0],
		ApplicationID: parts[1],
		UserID:        parts[2],
		UserRole:      parts[3],
	}, nil
}
