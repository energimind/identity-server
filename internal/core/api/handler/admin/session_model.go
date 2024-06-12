package admin

// session is a struct that contains session information.
type session struct {
	ID    string      `json:"id"`
	AppID string      `json:"appId"`
	User  sessionUser `json:"user"`
}

// sessionUser is a struct that contains sessionUser information.
type sessionUser struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}
