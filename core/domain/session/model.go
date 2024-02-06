package session

// UserInfo is a struct that contains user information.
// This info is returned by the identity provider upon successful authentication.
type UserInfo struct {
	ID         string
	Name       string
	GivenName  string
	FamilyName string
	Email      string
}

// Info is a struct that contains session information.
type Info struct {
	SessionID     string
	ApplicationID string
	UserInfo      UserInfo
}
