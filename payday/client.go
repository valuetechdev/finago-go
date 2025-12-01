//go:generate go tool -modfile=../go.tool.mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../api/openapi/payroll.json

package payday

import (
	"fmt"
	"net/http"
)

type PaydayClient struct {
	token      *Token
	httpClient *http.Client
	secret     string
	*ClientWithResponses
}

type Opts struct {
	Secret string
}

type Option func(*PaydayClient)

// WithHttpClient sets a custom http.Client. Defaults to [http.DefaultClient].
func WithHttpClient(client *http.Client) Option {
	return func(c *PaydayClient) {
		c.httpClient = client
	}
}

// Returns new [PaydayClient].
//
// You can reuse an already generated token and have it revalidated if it has
// expired, by using [PaydayClient.SetToken].
//
// You can provide options to customize the client behavior.
func New(secret string, options ...Option) *PaydayClient {
	client := &PaydayClient{httpClient: http.DefaultClient, secret: secret}
	for _, option := range options {
		option(client)
	}
	baseUrl := "https://payroll.24sevenoffice.com/api"
	c, err := NewClientWithResponses(
		baseUrl,
		WithRequestEditorFn(client.Intercept),
		WithHTTPClient(client.httpClient))

	if err != nil {
		panic(fmt.Errorf("failed to init client: %w", err))
	}
	client.ClientWithResponses = c
	return client
}
