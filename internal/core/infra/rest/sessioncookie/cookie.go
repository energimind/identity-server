package sessioncookie

import (
	"net/http"
)

func createCookie(r *http.Request, name, value, secret string) (*http.Cookie, error) {
	sc := getSecurityContext(r)

	encryptedValue, err := encryptCookie(value, secret)
	if err != nil {
		return nil, NewError("encrypt cookie error: %s", err)
	}

	const (
		days   = 60 * 60 * 24
		expiry = 30 * days
	)

	return &http.Cookie{
		Name:     name,
		Value:    encryptedValue,
		Path:     "/",
		Domain:   sc.domain,
		MaxAge:   expiry,
		Secure:   sc.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}, nil
}

func resetCookie(r *http.Request, name string) *http.Cookie {
	sc := getSecurityContext(r)

	const deleteCookieNow = -1

	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   sc.domain,
		MaxAge:   deleteCookieNow,
		Secure:   sc.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
