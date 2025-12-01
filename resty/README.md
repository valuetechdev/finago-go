[![go reference badge](https://pkg.go.dev/badge/github.com/valuetechdev/finago-go.svg)](https://pkg.go.dev/github.com/valuetechdev/finago-go/resty)

# Finago Office REST API Client

[Official docs](https://rest-api.developer.24sevenoffice.com/doc/v1/)

## Usage

```go
import "github.com/valuetechdev/finago-go/resty"

func yourFunc() (any, error) {
	client := New(&resty.Credentials{
		ClientId:       "your-client-id",
		ClientSecret:   "your-client-secret",
		OrganizationId: "your-org-id",
	})

	if err := client.Authenticate(); err != nil {
	  return nil, err
	}

	res, err := c.GetMeWithResponse(context.Background(), nil)
	if err != nil {
	  return nil, err
	}

	// Use the data
}
```
