package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

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

// Login completes the login process and returns the session information.
func (c *Client) Login(ctx context.Context, code, state string) (Session, error) {
	var result struct {
		SessionID     string `json:"sessionId"`
		ApplicationID string `json:"applicationId"`
		UserInfo      struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
			Email      string `json:"email"`
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
		return Session{}, newIdentityServerError("failed to complete login: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return Session{}, err
	}

	session := Session{
		SessionID:     result.SessionID,
		ApplicationID: result.ApplicationID,
		User: User{
			ID:         result.UserInfo.ID,
			Name:       result.UserInfo.Name,
			GivenName:  result.UserInfo.GivenName,
			FamilyName: result.UserInfo.FamilyName,
			Email:      result.UserInfo.Email,
		},
	}

	return session, nil
}

// Refresh refreshes the session and returns whether the session was refreshed.
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
		SetHeader("X-IS-SessionID", sessionID).
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
