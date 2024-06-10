package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

const sessionIDHeader = "X-IS-SessionID"

// Client is a client to interact with the identity service.
type Client struct {
	authEndpoint string
	rest         *resty.Client
}

// New returns a new instance of Client.
func New(authEndpoint string) *Client {
	const clientTimeout = 10 * time.Second

	return &Client{
		authEndpoint: authEndpoint + "/api/v1/auth",
		rest:         resty.New().SetTimeout(clientTimeout),
	}
}

// ProviderLink returns the link to the provider's login page.
func (c *Client) ProviderLink(ctx context.Context, appCode, providerCode, action string) (string, error) {
	var result struct {
		Link string `json:"link"`
	}

	rsp, err := c.rest.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"appCode":      appCode,
			"providerCode": providerCode,
			"action":       action,
		}).
		SetResult(&result).
		Get(c.authEndpoint + "/link")
	if err != nil {
		return "", newIdentityServerError("failed to get provider link: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return "", err
	}

	return result.Link, nil
}

// Login completes the login process and returns the session ID.
func (c *Client) Login(ctx context.Context, code, state string) (string, error) {
	var result struct {
		SessionID string `json:"sessionId"`
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
		return "", newIdentityServerError("failed to complete login: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return "", err
	}

	return result.SessionID, nil
}

// Session returns the session information.
func (c *Client) Session(ctx context.Context, sessionID string) (Session, error) {
	var result Session

	rsp, err := c.rest.R().
		SetContext(ctx).
		SetHeader(sessionIDHeader, sessionID).
		SetResult(&result).
		Get(c.authEndpoint + "/session")
	if err != nil {
		return Session{}, newIdentityServerError("failed to get session info: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return Session{}, err
	}

	return result, nil
}

// Refresh refreshes the session and returns whether the session was refreshed.
func (c *Client) Refresh(ctx context.Context, sessionID string) (bool, error) {
	var result struct {
		Refreshed bool `json:"refreshed"`
	}

	rsp, err := c.rest.R().
		SetContext(ctx).
		SetHeader(sessionIDHeader, sessionID).
		SetResult(&result).
		Put(c.authEndpoint + "/refresh")
	if err != nil {
		return false, newIdentityServerError("failed to refresh session: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return false, err
	}

	return result.Refreshed, nil
}

// Logout logs out the session.
func (c *Client) Logout(ctx context.Context, sessionID string) error {
	rsp, err := c.rest.R().
		SetContext(ctx).
		SetHeader(sessionIDHeader, sessionID).
		Delete(c.authEndpoint + "/logout")
	if err != nil {
		return newIdentityServerError("failed to logout: %v", err)
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

	return newIdentityServerError("%s (%d)", result.Error, rsp.StatusCode())
}
