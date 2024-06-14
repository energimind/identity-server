package admin

// session is a struct that contains session information.
type sessionInfo struct {
	ID      string      `json:"id"`
	RealmID string      `json:"realmId"`
	User    sessionUser `json:"user"`
}

// sessionUser is a struct that contains sessionUser information.
type sessionUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}
