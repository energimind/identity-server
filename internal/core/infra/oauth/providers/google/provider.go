// Package google implements the Google OAuth provider.
package google

import (
	"context"
	"net/http"

	"github.com/energimind/identity-server/internal/core/infra/oauth"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Provider is an OAuth provider for Google.
type Provider struct {
	config *oauth2.Config
}

// Ensure provider implements the oauth.Provider interface.
var _ oauth.Provider = (*Provider)(nil)

// NewProvider returns a new Google OAuth provider.
func NewProvider(config *oauth.Config) *Provider {
	return &Provider{
		config: providerConf(config),
	}
}

// GetAuthURL implements the oauth.Provider interface.
func (p *Provider) GetAuthURL(_ context.Context, state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Authorize implements the oauth.Provider interface.
func (p *Provider) Authorize(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, oauth.NewError("failed to exchange code for token: %v", err)
	}

	return token, nil
}

// RefreshAccessToken implements the oauth.Provider interface.
func (p *Provider) RefreshAccessToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	token, err := p.config.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, oauth.NewError("failed to refresh token: %v", err)
	}

	return token, nil
}

// RevokeAccessToken implements the oauth.Provider interface.
func (p *Provider) RevokeAccessToken(ctx context.Context, token *oauth2.Token) error {
	client := resty.New()

	rsp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetFormData(map[string]string{
			"token": token.AccessToken,
		}).
		Post("https://oauth2.googleapis.com/revoke")
	if err != nil {
		return oauth.NewError("failed to revoke token: %v", err)
	}

	if rsp.StatusCode() != http.StatusOK {
		return oauth.NewError("failed to revoke token: %s", rsp.Status())
	}

	return nil
}

// GetUserInfo implements the oauth.Provider interface.
func (p *Provider) GetUserInfo(ctx context.Context, token *oauth2.Token) (oauth.UserInfo, error) {
	client := resty.New()
	ui := map[string]any{}

	_, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetResult(&ui).
		Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return oauth.UserInfo{}, oauth.NewError("failed to get user info: %v", err)
	}

	str := func(key string) string {
		v := ui[key]

		if s, ok := v.(string); ok {
			return s
		}

		return ""
	}

	return oauth.UserInfo{
		ID:         str("sub"),
		BindID:     str("email"),
		Name:       str("name"),
		GivenName:  str("given_name"),
		FamilyName: str("family_name"),
		Email:      str("email"),
	}, nil
}

func providerConf(config *oauth.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
