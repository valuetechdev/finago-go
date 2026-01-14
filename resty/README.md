[![go reference badge](https://pkg.go.dev/badge/github.com/valuetechdev/finago-go.svg)](https://pkg.go.dev/github.com/valuetechdev/finago-go/resty)

# Finago Office REST API Client

[Official docs](https://rest-api.developer.24sevenoffice.com/doc/v1/)

## Usage

```go
import "github.com/valuetechdev/finago-go/resty"

func yourFunc() (any, error) {
	client := resty.New(&resty.Credentials{
		ClientId:       "your-client-id",
		ClientSecret:   "your-client-secret",
		OrganizationId: "your-org-id",
	})

	res, err := client.GetMeWithResponse(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return res.JSON200, nil
}
```

Authentication is handled automatically via OAuth2 client credentials. Tokens are refreshed when they expire.

### Reusing tokens

To persist and reuse tokens across sessions:

```go
client := resty.New(creds, resty.WithToken(existingToken))

// Get the current token to persist it
token, err := client.Token()
```

### Custom HTTP client

To use a custom HTTP client (e.g., with custom timeouts or transport):

```go
httpClient := &http.Client{Timeout: 30 * time.Second}
client := resty.New(creds, resty.WithBaseHTTPClient(httpClient))
```
