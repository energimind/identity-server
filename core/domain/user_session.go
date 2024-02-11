package domain

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
