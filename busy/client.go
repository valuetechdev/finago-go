//go:generate go tool -modfile=../go.tool.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../api/openapi/busy.json

package busy

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// ProdBaseURL is the production environment.
	ProdBaseURL = "https://api.busy.no"
	// DemoBaseURL is the demo/test environment. Get a demo workspace and a
	// token at https://demo.busy.no/demo/api.
	DemoBaseURL = "https://api.demo.busy.no"
)

// BusyClient is a client for the Finago Busy v2 REST API.
//
// Busy authenticates with a long-lived API token created by a workspace admin
// under the workspace's integration settings. The token is sent as a bearer
// token on every request.
type BusyClient struct {
	token   string
	baseURL string

	httpClient   *http.Client
	interceptors []RequestEditorFn
	retry        *retryTransport

	*ClientWithResponses
}

// Option configures a [BusyClient].
type Option func(*BusyClient)

// WithHttpClient sets a custom [http.Client]. Defaults to [http.DefaultClient].
func WithHttpClient(client *http.Client) Option {
	return func(c *BusyClient) {
		c.httpClient = client
	}
}

// WithURL overrides the API base URL. Defaults to [ProdBaseURL].
func WithURL(baseURL string) Option {
	return func(c *BusyClient) {
		c.baseURL = baseURL
	}
}

// WithDemo points the client at the demo environment, [DemoBaseURL].
func WithDemo() Option {
	return WithURL(DemoBaseURL)
}

// WithRequestInterceptor adds a request editor function that will be called
// before each request is sent.
func WithRequestInterceptor(fn RequestEditorFn) Option {
	return func(c *BusyClient) {
		c.interceptors = append(c.interceptors, fn)
	}
}

// New creates a new [BusyClient] with the given API token.
//
// Use [WithDemo] to talk to the demo environment while developing.
func New(token string, options ...Option) *BusyClient {
	client := &BusyClient{
		token:      token,
		baseURL:    ProdBaseURL,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	// Layer retrying over whatever transport the caller ended up with, so
	// WithRetry and WithHttpClient compose in either order.
	if client.retry != nil {
		httpClient := *client.httpClient
		client.retry.base = httpClient.Transport
		httpClient.Transport = client.retry
		client.httpClient = &httpClient
	}

	clientOptions := []ClientOption{
		WithHTTPClient(client.httpClient),
		WithRequestEditorFn(client.Intercept),
	}

	for _, interceptor := range client.interceptors {
		clientOptions = append(clientOptions, WithRequestEditorFn(interceptor))
	}

	c, err := NewClientWithResponses(client.baseURL, clientOptions...)
	if err != nil {
		panic(fmt.Errorf("failed to init client: %w", err))
	}
	client.ClientWithResponses = c
	return client
}

// Intercept sets the "Authorization: Bearer <token>" header on the request.
func (c *BusyClient) Intercept(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.token)

	return nil
}
