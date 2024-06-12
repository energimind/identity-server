package auth

// Session is a struct that contains session and user information.
type Session struct {
	Header Header
	User   User
}

// Header is a struct that contains session header information.
type Header struct {
	SessionID     string
	ApplicationID string
}

// User is a struct that contains user information.
type User struct {
	ID         string
	Name       string
	GivenName  string
	FamilyName string
	Email      string
}
