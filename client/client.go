package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is a client to interact with the identity service.
type Client struct {
	baseURL string
	rest    *resty.Client
}

// New returns a new instance of Client.
// The baseURL is the base URL of the identity service.
func New(baseURL string) *Client {
	const clientTimeout = 10 * time.Second

	return &Client{
		baseURL: baseURL + "/api/v1/sessions/",
		rest:    resty.New().SetTimeout(clientTimeout),
	}
}

// Session returns the session information.
func (c *Client) Session(ctx context.Context, sessionID string) (Session, error) {
	var result Session

	rsp, err := c.newRequest(ctx).SetResult(&result).Get(c.baseURL + sessionID)
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

	rsp, err := c.newRequest(ctx).SetResult(&result).Put(c.baseURL + sessionID + "/refresh")
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
	rsp, err := c.newRequest(ctx).Delete(c.baseURL + sessionID)
	if err != nil {
		return newIdentityServerError("failed to logout: %v", err)
	}

	return processErrorResponse(rsp)
}

func (c *Client) newRequest(ctx context.Context) *resty.Request {
	return c.rest.R().SetContext(ctx)
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
