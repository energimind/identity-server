package sessioncookie

import (
	"net/url"

	"github.com/energimind/identity-server/core/domain"
	"github.com/gin-gonic/gin"
)

// Provider defines methods for creating, resetting cookies and parsing cookies.
//
// It is used to create and parse user sessions stored in cookies.
type Provider struct {
	cookieName string
	secret     string
}

// NewProvider returns a new instance of Provider.
func NewProvider(cookieName, secret string) *Provider {
	return &Provider{
		cookieName: cookieName,
		secret:     padSecret(secret),
	}
}

// CreateCookie creates a cookie with the given user session.
func (p *Provider) CreateCookie(c *gin.Context, us domain.UserSession) error {
	cookie, err := createCookie(c.Request, p.cookieName, serializeUserSession(us), p.secret)
	if err != nil {
		return err
	}

	c.SetSameSite(cookie.SameSite)
	c.SetCookie(p.cookieName, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

	return nil
}

// ResetCookie resets the cookie.
func (p *Provider) ResetCookie(c *gin.Context) error {
	cookie := resetCookie(c.Request, p.cookieName)

	c.SetSameSite(cookie.SameSite)
	c.SetCookie(p.cookieName, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

	return nil
}

// ParseCookie parses the cookie and returns the user session.
func (p *Provider) ParseCookie(c *gin.Context) (domain.UserSession, error) {
	cookie, err := c.Request.Cookie(p.cookieName)
	if err != nil {
		return domain.UserSession{}, NewError("cookie not found: %s", err)
	}

	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return domain.UserSession{}, NewError("cookie decode error: %s", err)
	}

	decryptedValue, err := decryptCookie(decodedValue, p.secret)
	if err != nil {
		return domain.UserSession{}, NewError("decrypt cookie error: %s", err)
	}

	us, err := deserializeUserSession(decryptedValue)
	if err != nil {
		return domain.UserSession{}, err
	}

	return us, nil
}
