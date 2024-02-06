package cookie

import (
	"net/http"

	"github.com/energimind/identity-service/core/domain/auth"
)

// Provider is a cookie provider.
//
// It provides methods for getting, resetting and verifying cookies.
//
// It implements the auth.CookieProvider and auth.CookieVerifier interfaces.
type Provider struct {
	secret string
}

// NewProvider creates a new cookie provider.
func NewProvider(secret string) *Provider {
	return &Provider{
		secret: padSecret(secret),
	}
}

// Ensure that Provider implements the CookieProvider interface.
var _ auth.CookieProvider = (*Provider)(nil)

// Ensure that Provider implements the CookieVerifier interface.
var _ auth.CookieVerifier = (*Provider)(nil)

// CreateCookie creates a cookie with the given name and value.
func (p *Provider) CreateCookie(r *http.Request, name, value string) (*http.Cookie, error) {
	sc, err := getSecurityContext(r)
	if err != nil {
		return nil, NewError("get security context error: %s", err)
	}

	encryptedValue, err := encryptCookie(value, p.secret)
	if err != nil {
		return nil, NewError("encrypt cookie error: %s", err)
	}

	const neverExpires = 0

	return &http.Cookie{
		Name:     name,
		Value:    encryptedValue,
		Path:     "/",
		Domain:   sc.domain,
		MaxAge:   neverExpires,
		Secure:   sc.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}, nil
}

// ResetCookie resets a cookie.
func (p *Provider) ResetCookie(r *http.Request, name string) (*http.Cookie, error) {
	sc, err := getSecurityContext(r)
	if err != nil {
		return nil, NewError("get security context error: %s", err)
	}

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
	}, nil
}

// VerifyCookie verifies a cookie.
func (p *Provider) VerifyCookie(r *http.Request, cookie *http.Cookie) (string, error) {
	sc, err := getSecurityContext(r)
	if err != nil {
		return "", NewError("get security context error: %s", err)
	}

	if cookie.Domain != sc.domain {
		return "", NewError("invalid cookie domain: %s", cookie.Domain)
	}

	if cookie.Secure != sc.secure {
		return "", NewError("invalid cookie secure flag: %t", cookie.Secure)
	}

	decryptedValue, err := decryptCookie(cookie.Value, p.secret)
	if err != nil {
		return "", NewError("decrypt cookie error: %s", err)
	}

	return decryptedValue, nil
}
