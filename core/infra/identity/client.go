package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/energimind/identity-service/core/domain"
	"github.com/energimind/identity-service/core/domain/admin"
	"github.com/go-resty/resty/v2"
)

// adminActor is the actor for the admin role.
//
//nolint:gochecknoglobals // it is a constant
var adminActor = admin.Actor{Role: admin.SystemRoleAdmin}

// Client is a client to interact with the identity service.
type Client struct {
	authEndpoint string
	userFinder   admin.UserFinder
	rest         *resty.Client
}

// Ensure Client implements admin.IdentityClient.
var _ admin.IdentityClient = (*Client)(nil)

// NewClient returns a new instance of Client.
func NewClient(authEndpoint string, userFinder admin.UserFinder) *Client {
	const clientTimeout = 10 * time.Second

	return &Client{
		authEndpoint: authEndpoint + "/api/v1/auth",
		userFinder:   userFinder,
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
func (c *Client) Login(ctx context.Context, code, state string) (admin.Session, admin.User, error) {
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
		return admin.Session{}, admin.User{}, domain.NewGatewayError("failed to complete login: %v", err)
	}

	if err := processErrorResponse(rsp); err != nil {
		return admin.Session{}, admin.User{}, err
	}

	user, err := c.userFinder.GetUserByEmail(
		ctx,
		adminActor,
		admin.ID(result.ApplicationID),
		result.UserInfo.Email,
	)
	if err != nil {
		//nolint:wrapcheck // this is already a domain error
		return admin.Session{}, admin.User{}, err
	}

	session := admin.Session{
		SessionID:     result.SessionID,
		ApplicationID: result.ApplicationID,
	}

	return session, user, nil
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
