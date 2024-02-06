package auth

import (
	"net/http"
)

// CookieProvider defines methods for creating and resetting cookies.
type CookieProvider interface {
	CreateCookie(r *http.Request, name, value string) (*http.Cookie, error)
	ResetCookie(r *http.Request, name string) (*http.Cookie, error)
}

// CookieVerifier defines a method for verifying cookies.
type CookieVerifier interface {
	VerifyCookie(r *http.Request, cookie *http.Cookie) (string, error)
}
