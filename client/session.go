package client

// User is a struct that contains user information.
type User struct {
	ID         string
	Name       string
	GivenName  string
	FamilyName string
	Email      string
}

// Session is a struct that contains session information.
type Session struct {
	SessionID     string
	ApplicationID string
	User          User
}
