//go:generate go tool -modfile=../go.tool.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../api/openapi/rest.yaml

package resty

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// RestyClient is a client for the 24SevenOffice REST API.
// It handles OAuth2 authentication automatically and refreshes tokens when they expire.
type RestyClient struct {
	tokenSource oauth2.TokenSource

	interceptors []RequestEditorFn

	initialToken   *oauth2.Token
	baseHTTPClient *http.Client

	*ClientWithResponses
}

// Credentials contains the OAuth2 client credentials for authenticating
// with the 24SevenOffice API.
type Credentials struct {
	ClientId       string
	ClientSecret   string
	OrganizationId string
}

// Option configures a [RestyClient].
type Option func(*RestyClient)

// WithToken sets an existing token to be reused. The token will be
// automatically refreshed when it expires.
func WithToken(token *oauth2.Token) Option {
	return func(c *RestyClient) {
		c.initialToken = token
	}
}

// WithRequestInterceptor adds a request editor function that will be called
// before each request is sent.
func WithRequestInterceptor(fn RequestEditorFn) Option {
	return func(c *RestyClient) {
		c.interceptors = append(c.interceptors, fn)
	}
}

// WithHttpClient sets a custom HTTP client to use as the base for OAuth2 requests.
// The provided client's transport will be wrapped with OAuth2 authentication.
func WithHttpClient(client *http.Client) Option {
	return func(c *RestyClient) {
		c.baseHTTPClient = client
	}
}

// New creates a new [RestyClient] with the given credentials.
// Authentication is handled automatically via OAuth2 client credentials.
// Tokens are refreshed when they expire.
//
// Use [WithToken] to reuse an existing token across sessions.
func New(credentials *Credentials, options ...Option) *RestyClient {
	conf := &clientcredentials.Config{
		ClientID:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		TokenURL:     "https://login.24sevenoffice.com/oauth/token",
		EndpointParams: map[string][]string{
			"login_organization": {credentials.OrganizationId},
			"audience":           {"https://api.24sevenoffice.com"},
		},
	}

	baseUrl := "https://rest.api.24sevenoffice.com/v1"
	client := &RestyClient{}

	for _, option := range options {
		option(client)
	}

	// Use custom base HTTP client if provided
	ctx := context.Background()
	if client.baseHTTPClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client.baseHTTPClient)
	}

	// Create a token source that reuses the initial token (if provided)
	// and automatically refreshes using client credentials when expired
	client.tokenSource = oauth2.ReuseTokenSource(
		client.initialToken,
		conf.TokenSource(ctx),
	)

	// Create an HTTP client that automatically handles token auth
	httpClient := oauth2.NewClient(ctx, client.tokenSource)

	clientOptions := []ClientOption{
		WithHTTPClient(httpClient),
	}

	for _, interceptor := range client.interceptors {
		clientOptions = append(clientOptions, WithRequestEditorFn(interceptor))
	}

	c, err := NewClientWithResponses(
		baseUrl,
		clientOptions...,
	)
	if err != nil {
		panic(fmt.Errorf("failed to init client: %w", err))
	}
	client.ClientWithResponses = c
	return client
}

// Token returns the current OAuth2 token. If the token has expired,
// it will be automatically refreshed.
func (c *RestyClient) Token() (*oauth2.Token, error) {
	return c.tokenSource.Token()
}
