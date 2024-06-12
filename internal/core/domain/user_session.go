package domain

// UserSession is a struct that contains user session information.
type UserSession struct {
	SessionID string
	RealmID   string
	UserID    string
	UserRole  string
}

// NewUserSession creates a new UserSession with the given parameters.
func NewUserSession(sessionID, realmID, userID, userRole string) UserSession {
	return UserSession{
		SessionID: sessionID,
		RealmID:   realmID,
		UserID:    userID,
		UserRole:  userRole,
	}
}
