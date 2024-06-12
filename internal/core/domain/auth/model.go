package auth

// Session is a struct that contains session and user information.
type Session struct {
	Header Header `json:"header"`
	User   User   `json:"user"`
}

// Header is a struct that contains session header information.
type Header struct {
	SessionID     string `json:"sessionId"`
	ApplicationID string `json:"applicationId"`
}

// User is a struct that contains user information.
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
	Email      string `json:"email"`
}
