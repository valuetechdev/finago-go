[![go reference badge](https://pkg.go.dev/badge/github.com/valuetechdev/finago-go.svg)](https://pkg.go.dev/github.com/valuetechdev/finago-go)

# Finago API clients for Go

This package contains a generated API clients for Finago's APIs.

- Finago Office SOAP : [`soapy`](soapy/README.md) ([API](https://developer.24sevenoffice.com/docs/))
- Finago Office REST: [`resty`](resty/README.md) ([API](https://rest-api.developer.24sevenoffice.com/doc/v1/))
- Finago Payday: [`payday`](payday/README.md) ([API](https://swagger.api.24sevenoffice.com/?url=https://me.24sevenoffice.com/swagger.json))

## Usage

```bash
go get github.com/valuetechdev/finago-go
```

## Development

Tooling is managed by [mise](https://mise.jdx.dev), secrets by
[fnox](https://fnox.jdx.dev/) and linting by [hk](https://hk.jdx.dev).

```bash
mise install # installs tools, verifies secrets and installs git hooks
mise tasks   # lists available tasks
```

| Task                 | Description                                   |
| -------------------- | --------------------------------------------- |
| `mise run api`       | Fetch the latest OpenAPI specifications       |
| `mise run generate`  | Regenerate the API clients                    |
| `mise run test`      | Run the offline tests, no secrets needed      |
| `mise run test:full` | Run every test against the live APIs          |
| `mise run check`     | Lint unstaged code (`check:all`, `check:fix`) |
| `mise run tidy`      | `go mod tidy` and `go fmt`                    |
| `mise run bump`      | Bump to the next version                      |

### Tests

`mise run test` runs `go test -short ./...`, which skips everything that talks
to a Finago API. It needs no credentials and is what CI and the git hooks run.

`mise run test:full` hits the live APIs and needs credentials for a Finago
tenant. It runs the suite through `fnox exec`, so the following environment
variables must resolve:

| Variable                  | Used by  | What it is                   |
| ------------------------- | -------- | ---------------------------- |
| `TFSO_SOAP_APPLICATIONID` | `soapy`  | SOAP application ID (a GUID) |
| `TFSO_SOAP_USERNAME`      | `soapy`  | SOAP user, an e-mail address |
| `TFSO_SOAP_PASSWORD`      | `soapy`  | SOAP user password           |
| `TFSO_REST_APP_ID`        | `resty`  | REST application ID          |
| `TFSO_REST_SECRET`        | `resty`  | REST client secret           |
| `TFSO_PAYROLL_SECRET`     | `payday` | Payday API credential        |

The `resty` tests are pinned to the demo organization `543819716587312`, so the
REST credentials need access to it.

### Configuring secrets

Use [fnox](https://fnox.jdx.dev/), and set up a `fnox.loca.toml`. If you use
1Password, it can look like this:

```toml
default_provider = "onepass"

[providers]
onepass = { type = "1password", vault = "<value name>", account = "<account name>.1password.com" }

[secrets]
TFSO_SOAP_APPLICATIONID = { provider = "onepass", value = "<secret reference>" }
TFSO_SOAP_USERNAME = { provider = "onepass", value = "<secret reference>" }
TFSO_SOAP_PASSWORD = { provider = "onepass", value = "<secret reference>" }

TFSO_REST_APP_ID = { provider = "onepass", value = "<secret reference>" }
TFSO_REST_SECRET = { provider = "onepass", value = "<secret reference>" }

TFSO_PAYROLL_SECRET = { provider = "onepass", value = "<secret reference>" }
```

Otherwise point the same six keys at your own tenant with whichever [fnox
provider](https://fnox.jdx.dev/reference/providers) you prefer — `fnox init`
walks you through picking one, and `fnox set TFSO_REST_SECRET` stores a value.
A provider is mandatory: `fnox check` fails with "No providers configured" if
`fnox.local.toml` only lists secrets.

Verify the setup with:

```bash
fnox check           # config resolves and every secret is reachable
fnox exec -- env | grep TFSO_
```
