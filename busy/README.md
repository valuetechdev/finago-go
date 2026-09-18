[![go reference badge](https://pkg.go.dev/badge/github.com/valuetechdev/finago-go.svg)](https://pkg.go.dev/github.com/valuetechdev/finago-go/busy)

# Finago Busy API Client

[Official docs](https://api.busy.no/v2/)

## Usage

```go
import "github.com/valuetechdev/finago-go/busy"

func yourFunc() (any, error) {
	client := busy.New("your-api-token")

	res, err := client.GetAllUsersWithResponse(context.Background(), &busy.GetAllUsersParams{})
	if err != nil {
		return nil, err
	}

	return res.JSON200, nil
}
```

The API token is created by a workspace admin under the workspace's integration
settings, and is sent as `Authorization: Bearer <token>` on every request. A
token is valid for at most a year, so it has to be rotated.

### Demo environment

Busy has a separate demo environment, and **a token is only valid in the
environment it was created in** — a production token gets a 401 from the demo
host and vice versa. Grab a demo workspace and token from
<https://demo.busy.no/demo/api>, then:

```go
client := busy.New("your-demo-token", busy.WithDemo())
```

`busy.WithURL` points the client at any other host.

### Rate limiting

Busy throttles with a **429** and tells you how long to wait. `busy.WithRetry`
retries those requests for you, reading the wait from `Retry-After`, then
`RateLimit-Reset`, then falling back to exponential backoff:

```go
client := busy.New(token, busy.WithRetry())
```

It is off by default, since it makes a call block for as long as the limiter
asks. Only 429 is retried — other failures, 5xx included, come back untouched —
and the wait is abandoned if the request's context is cancelled. Tune it with
`busy.WithRetryAttempts` (default 3) and `busy.WithRetryMaxWait` (default 60s);
a throttling window longer than the cap is handed back as a 429 rather than
slept through.

`WithRetry` wraps the transport of the client from `WithHttpClient`, so the two
compose in either order, and neither your client nor `http.DefaultClient` is
modified.

### Custom HTTP client

```go
httpClient := &http.Client{Timeout: 30 * time.Second}
client := busy.New(token, busy.WithHttpClient(httpClient))
```

## Things to know

- The schema is OpenAPI 3.1 and is generated as published, apart from
  `overlay.yaml`, which maps the `format: email` fields onto this package's own
  `busy.Email` type. Busy returns `""` for users without an e-mail address, and
  the `openapi_types.Email` those fields would otherwise generate rejects that,
  failing the decode of the entire response. `busy.Email` validates every
  non-empty value the same way and accepts `""`.
- The paths in the schema include the `/v2` prefix, so the base URL is the bare
  host (`https://api.busy.no`).
- The API is rate limited. Responses carry `RateLimit-Limit`,
  `RateLimit-Remaining` and `RateLimit-Reset`; a **429** also carries
  `Retry-After`. `busy.WithRetry` handles the backoff; without it, back off for
  at least `RateLimit-Reset` seconds when throttled.
- List endpoints are paginated with `limit`/`offset`. For syncs, order ascending
  by `updatedAt` and filter on `updatedFrom`.
