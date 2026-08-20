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
[fnox](https://github.com/jdx/fnox) (1Password) and linting by
[hk](https://hk.jdx.dev).

```bash
mise install # installs tools, verifies secrets and installs git hooks
mise tasks   # lists available tasks
```

| Task              | Description                                          |
| ----------------- | ---------------------------------------------------- |
| `mise run api`    | Fetch the latest OpenAPI specifications               |
| `mise run generate` | Regenerate the API clients                         |
| `mise run test`   | Run tests with secrets injected from 1Password        |
| `mise run check`  | Lint unstaged code (`check:all`, `check:fix`)         |
| `mise run tidy`   | `go mod tidy` and `go fmt`                            |
| `mise run bump`   | Bump to the next version                              |
