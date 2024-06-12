package client

// Session is a struct that contains session and user information.
type Session struct {
	Header Header `json:"header"`
	User   User   `json:"user"`
}

// Header is a struct that contains session header information.
type Header struct {
	SessionID string `json:"sessionId"`
	RealmID   string `json:"realmId"`
}

// User is a struct that contains user information.
type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}
