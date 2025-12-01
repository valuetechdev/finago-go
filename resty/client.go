//go:generate go tool -modfile=../go.tool.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../api/openapi/rest.yaml

package resty

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type RestyClient struct {
	token *oauth2.Token
	conf  *clientcredentials.Config

	interceptors []RequestEditorFn

	httpClient *http.Client
	*ClientWithResponses
}

type Credentials struct {
	ClientId       string
	ClientSecret   string
	OrganizationId string
}

type Option func(*RestyClient)

// WithHttpClient sets a custom http.Client. Defaults to [http.DefaultClient].
func WithHttpClient(client *http.Client) Option {
	return func(c *RestyClient) {
		c.httpClient = client
	}
}

func WithRequestInterceptor(fn RequestEditorFn) Option {
	return func(c *RestyClient) {
		c.interceptors = append(c.interceptors, fn)
	}
}

// Returns new [RestyClient].
//
// You can reuse an already generated token and have it revalidated if it has
// expired, by using [RestyClient.SetToken].
//
// You can provide options to customize the client behavior.
//
// WARNING: The client hasn't been tested.
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
	client := &RestyClient{conf: conf, httpClient: http.DefaultClient}

	for _, option := range options {
		option(client)
	}

	clientOptions := []ClientOption{
		WithRequestEditorFn(client.InterceptToken),
		WithHTTPClient(client.httpClient),
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
