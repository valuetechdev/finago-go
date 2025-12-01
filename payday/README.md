[![go reference badge](https://pkg.go.dev/badge/github.com/valuetechdev/finago-go.svg)](https://pkg.go.dev/github.com/valuetechdev/finago-go/payday)

# Finago Payday API Client

[Official docs](https://swagger.api.24sevenoffice.com/?url=https://me.24sevenoffice.com/swagger.json)

## Usage

```go
import "github.com/valuetechdev/finago-go/payday"

func yourFunc() (any, error) {
   c := payday.New("your-secret")
   if err := c.Authenticate(); err != nil {
      return nil, err
   }

	 absenceRes, err := c.GetAbsenceV2WithResponse(context.TODO())
   if err != nil {
      return nil, err
   }

   // Use the data
}
```

## Things to know

- The original schema is using OpenAPI/Swagger 2.0.
- The schema has been altered to with defined schemas/models that was missing.
