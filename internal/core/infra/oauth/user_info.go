package oauth

// UserInfo represents the user information returned by the OAuth provider.
type UserInfo struct {
	ID         string
	Name       string
	GivenName  string
	FamilyName string
	Email      string
}
