package client

// info is a struct that contains session and user information.
type info struct {
	SessionInfo sessionInfo `json:"sessionInfo"`
	UserInfo    userInfo    `json:"userInfo"`
}

// sessionInfo is a struct that contains session information.
type sessionInfo struct {
	SessionID     string `json:"sessionId"`
	ApplicationID string `json:"applicationId"`
}

// userInfo is a struct that contains user information.
type userInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
}
