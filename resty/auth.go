package resty

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type TokenOpts struct {
	ClientId       string
	ClientSecret   string
	OrganizationId string
}

// Authenticates the current c
func (c *RestyClient) Authenticate() error {
	token, err := c.conf.Token(context.Background())
	if err != nil {
		return err
	}

	c.token = token
	return nil
}

// Returns the oauth2 token
func (c *RestyClient) GetToken() *oauth2.Token {
	return c.token
}

// Set the oauth2 token
func (c *RestyClient) SetToken(token *oauth2.Token) {
	c.token = token
}

// Returns wether the token is valid or not
func (c *RestyClient) IsTokenValid() bool {
	return c.token != nil && c.token.Valid()
}

// Intercept the req with an "Authorization: bearer <accessToken>" header.
func (c *RestyClient) InterceptToken(ctx context.Context, req *http.Request) error {
	if !c.IsTokenValid() {
		return fmt.Errorf("token is not valid")
	}
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	return nil
}
