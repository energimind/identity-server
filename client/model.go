package client

// Session is a struct that contains session and user information.
type Session struct {
	SessionInfo SessionInfo `json:"sessionInfo"`
	UserInfo    UserInfo    `json:"userInfo"`
}

// SessionInfo is a struct that contains session information.
type SessionInfo struct {
	SessionID     string `json:"sessionId"`
	ApplicationID string `json:"applicationId"`
}

// UserInfo is a struct that contains user information.
type UserInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
}
