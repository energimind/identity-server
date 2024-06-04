package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/energimind/identity-server/core/domain"
	"github.com/energimind/identity-server/core/domain/admin"
	"github.com/go-resty/resty/v2"
)

// Client is a client to interact with the identity service.
type Client struct {
	authEndpoint string
	rest         *resty.Client
}

// Ensure Client implements admin.IdentityClient.
var _ admin.IdentityClient = (*Client)(nil)

// NewClient returns a new instance of Client.
func NewClient(authEndpoint string) *Client {
	const clientTimeout = 10 * time.Second

	return &Client{
		authEndpoint: authEndpoint + "/api/v1/auth",
		rest:         resty.New().SetTimeout(clientTimeout),
	}
}

// ProviderLink implements admin.IdentityClient.
func (c *Client) ProviderLink(ctx context.Context, appCode, providerCode string) (string, error) {
	var result struct {
		Link string `json:"link"`
	}

	rsp, err := c.rest.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"appCode":      appCode,
			"providerCode": providerCode,
		}).
		SetResult(&result).
		Get(c.authEndpoint + "/link")
	if err != nil {
		return "", domain.NewGatewayError("failed to get provider link: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return "", err
	}

	return result.Link, nil
}

// Login implements admin.IdentityClient.
func (c *Client) Login(ctx context.Context, code, state string) (admin.Session, error) {
	var result struct {
		SessionID     string `json:"sessionId"`
		ApplicationID string `json:"applicationId"`
		UserInfo      struct {
			Email string `json:"email"`
		} `json:"userInfo"`
	}

	rsp, err := c.rest.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"code":  code,
			"state": state,
		}).
		SetResult(&result).
		Post(c.authEndpoint + "/login")
	if err != nil {
		return admin.Session{}, domain.NewGatewayError("failed to complete login: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return admin.Session{}, err
	}

	session := admin.Session{
		SessionID:     result.SessionID,
		ApplicationID: result.ApplicationID,
		UserEmail:     result.UserInfo.Email,
	}

	return session, nil
}

// Refresh implements admin.IdentityClient.
func (c *Client) Refresh(ctx context.Context, sessionID string) (bool, error) {
	var result struct {
		Refreshed bool `json:"refreshed"`
	}

	rsp, err := c.rest.R().
		SetContext(ctx).
		SetHeader("X-IS-SessionID", sessionID).
		SetResult(&result).
		Put(c.authEndpoint + "/refresh")
	if err != nil {
		return false, domain.NewGatewayError("failed to refresh session: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return false, err
	}

	return result.Refreshed, nil
}

// Logout implements admin.IdentityClient.
func (c *Client) Logout(ctx context.Context, sessionID string) error {
	rsp, err := c.rest.R().
		SetContext(ctx).
		SetHeader("X-IS-SessionID", sessionID).
		Delete(c.authEndpoint + "/logout")
	if err != nil {
		return domain.NewGatewayError("failed to logout: %v", err)
	}

	return processErrorResponse(rsp)
}

func processErrorResponse(rsp *resty.Response) error {
	if rsp.IsSuccess() {
		return nil
	}

	var result struct {
		Error string `json:"error"`
	}

	// ignore error if the response is not JSON
	_ = json.Unmarshal(rsp.Body(), &result)

	if result.Error == "" {
		result.Error = "unspecified"
	}

	return domain.NewGatewayError("%s (%d)", result.Error, rsp.StatusCode())
}
