package login

import (
	"context"
	"net/http"

	"github.com/energimind/identity-service/core/domain"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getAuthCodeURL(_ context.Context, config *config, state string) string {
	return providerConf(config).AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func exchangeCodeForAccessToken(ctx context.Context, config *config, code string) (*oauth2.Token, error) {
	token, err := providerConf(config).Exchange(ctx, code)
	if err != nil {
		return nil, domain.NewAccessDeniedError("failed to exchange code for token: %v", err)
	}

	return token, nil
}

func refreshAccessToken(ctx context.Context, config *config, token *oauth2.Token) (*oauth2.Token, error) {
	token, err := providerConf(config).TokenSource(ctx, token).Token()
	if err != nil {
		return nil, domain.NewAccessDeniedError("failed to refresh token: %v", err)
	}

	return token, nil
}

func revokeAccessToken(ctx context.Context, token *oauth2.Token) error {
	client := resty.New()

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token.AccessToken).
		SetFormData(map[string]string{
			"token": token.AccessToken,
		}).
		Post("https://oauth2.googleapis.com/revoke")
	if err != nil {
		return domain.NewAccessDeniedError("failed to revoke token: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return domain.NewAccessDeniedError("failed to revoke token: %s", resp.Status())
	}

	return nil
}

func providerConf(config *config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
