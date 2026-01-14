package resty

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"
	u "github.com/valuetechdev/finago-go/resty/internal"
	"golang.org/x/oauth2"
)

// SO24 project: "WooCommerce-integrasjon (DEMO)"
const orgId = "543819716587312"

var cachedClient *RestyClient

func getClient() *RestyClient {
	if cachedClient != nil {
		return cachedClient
	}
	cachedClient = New(&Credentials{
		ClientId:       os.Getenv("TFSO_REST_APP_ID"),
		ClientSecret:   os.Getenv("TFSO_REST_SECRET"),
		OrganizationId: orgId,
	})
	return cachedClient
}

func TestClientInitialization(t *testing.T) {
	require := require.New(t)

	c := New(&Credentials{
		ClientId:       os.Getenv("TFSO_REST_APP_ID"),
		ClientSecret:   os.Getenv("TFSO_REST_SECRET"),
		OrganizationId: orgId,
	})

	// Token is fetched automatically on first request, but we can also get it directly
	token, err := c.Token()
	require.NoError(err, "should fetch token")
	require.NotNil(token, "token should not be nil")
	require.NotEmpty(token.AccessToken, "access token should not be empty")
	require.Equal("Bearer", token.TokenType, "token type should be Bearer")

	t.Logf("Authentication successful. Token expires: %v", token.Expiry)

	res, err := c.GetAccountsWithResponse(context.Background(), &GetAccountsParams{})
	require.NoError(err, "GetAccounts should not error")
	require.NotNil(res, "GetAccounts should not be nil")
	require.Equal(http.StatusOK, res.StatusCode(), "GetAccounts status should be OK", string(res.Body))
	require.NotNil(res.JSON200, "GetAccounts JSON200 should not be nil")
}

func TestClientTokenManagement(t *testing.T) {
	require := require.New(t)

	testToken := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}

	c := New(&Credentials{
		ClientId:       "test-client",
		ClientSecret:   "test-secret",
		OrganizationId: "test-org",
	}, WithToken(testToken))

	token, err := c.Token()
	require.NoError(err)
	require.Equal(testToken.AccessToken, token.AccessToken, "access token should match")
}

func TestCreatePrivateCustomer(t *testing.T) {
	require := require.New(t)

	c := getClient()

	cPostRequest := CustomerPostRequest{
		Email: &EmailsDto{
			Billing: u.R("billing@example.org"),
			Contact: u.R("contact@example.org"),
		},
		Phone: u.R("+47-12345678"),

		Address: &AddressesDto{
			Visit: &VisitAddress{
				Street:             u.R("Torgallmenningen 1"),
				PostalArea:         u.R("Bergen"),
				PostalCode:         u.R("5009"),
				CountryCode:        u.R("NO"),
				CountrySubdivision: u.R("Vestland"),
			},
		},
		IsSupplier: u.R(false),
	}
	err := cPostRequest.FromCustomerPostRequest1(CustomerPostRequest1{
		IsCompany: CustomerPostRequest1IsCompany(false),
		Person: FirstnameLastnameDto{
			FirstName: u.R("John"),
			LastName:  u.R("Doe"),
		},
	})
	require.NoError(err)

	res, err := c.CreateCustomerWithResponse(t.Context(), cPostRequest)
	require.NoError(err)
	require.NotNil(res.JSON200, "no private customer was created")
}

func TestCreateCompanyCustomer(t *testing.T) {
	require := require.New(t)

	c := getClient()

	cPostRequest := CustomerPostRequest{
		Email: &EmailsDto{
			Billing: u.R("billing@example.org"),
			Contact: u.R("contact@example.org"),
		},
		Phone: u.R("+47-12345678"),
		Address: &AddressesDto{
			Visit: &VisitAddress{
				Street:             u.R("Torgallmenningen 1"),
				PostalArea:         u.R("Bergen"),
				PostalCode:         u.R("5009"),
				CountryCode:        u.R("NO"),
				CountrySubdivision: u.R("Vestland"),
			},
		},
		IsSupplier:         u.R(false),
		OrganizationNumber: u.R("123456789"),
	}
	err := cPostRequest.FromCustomerPostRequest0(CustomerPostRequest0{
		IsCompany: CustomerPostRequest0IsCompany(true),
		Name:      "ACME INC.",
	})
	require.NoError(err)

	res, err := c.CreateCustomerWithResponse(t.Context(), cPostRequest)
	require.NoError(err)
	require.NotNil(res.JSON200, "no company customer was created")
}

func TestRetrieveProductUnits(t *testing.T) {
	require := require.New(t)
	c := getClient()
	res, err := c.GetUnitsWithResponse(t.Context())
	require.NoError(err)
	require.NotNil(res.JSON200, "did not fetch product units")
}

func TestCreateProduct(t *testing.T) {
	require := require.New(t)

	c := getClient()

	cPostRequest := CustomerPostRequest{
		Email: &EmailsDto{
			Billing: u.R("billing@example.org"),
			Contact: u.R("contact@example.org"),
		},
		Phone: u.R("+47-12345678"),
		Address: &AddressesDto{
			Visit: &VisitAddress{
				Street:             u.R("Torgallmenningen 1"),
				PostalArea:         u.R("Bergen"),
				PostalCode:         u.R("5009"),
				CountryCode:        u.R("NO"),
				CountrySubdivision: u.R("Vestland"),
			},
		},
		IsSupplier:         u.R(true),
		OrganizationNumber: u.R("123456789"),
	}
	err := cPostRequest.FromCustomerPostRequest0(CustomerPostRequest0{
		IsCompany: CustomerPostRequest0IsCompany(true),
		Name:      "ACME INC.",
	})
	require.NoError(err)
	resCustomer, err := c.CreateCustomerWithResponse(t.Context(), cPostRequest)
	require.NoError(err)
	require.NotNil(resCustomer.JSON200, "no customer was created")

	pPostRequest := ProductRequestPost{
		Name:         u.R("Badeball"),
		Number:       u.R(u.RandSeq(10)),
		Type:         u.R(Default),
		Status:       u.R(ProductStatusEnumActive),
		Description:  u.R("Rund og flytende"),
		CostPrice:    u.R(float32(30)),
		IndirectCost: u.R(float32(20)),
		SalesPrice:   u.R(float32(100)),
		Stock: &StockDto{
			IsManaged: u.R(true),
			Quantity:  u.R(float32(129)),
			Location:  u.R("B-301"),
		},
		Category: &CategoryRequest{Id: u.R(-1)},
	}
	resProduct, err := c.CreateProductWithResponse(t.Context(), pPostRequest)
	require.NoError(err)
	require.NotNil(resProduct.JSON201, "no product was created")
}

func TestCreateOrder(t *testing.T) {
	require := require.New(t)

	c := getClient()

	pPostRequest := ProductRequestPost{
		Name:         u.R("Badeball"),
		Number:       u.R(u.RandSeq(10)),
		Type:         u.R(Default),
		Status:       u.R(ProductStatusEnumActive),
		Description:  u.R("Rund og flytende"),
		CostPrice:    u.R(float32(30)),
		IndirectCost: u.R(float32(20)),
		SalesPrice:   u.R(float32(100)),
		Stock: &StockDto{
			IsManaged: u.R(true),
			Quantity:  u.R(float32(129)),
			Location:  u.R("B-301"),
		},
		Category: &CategoryRequest{Id: u.R(-1)},
	}
	resProduct, err := c.CreateProductWithResponse(t.Context(), pPostRequest)
	require.NoError(err)
	require.NotNil(resProduct.JSON201, "no product was created")

	cPostRequest := CustomerPostRequest{
		Email: &EmailsDto{
			Billing: u.R("billing@example.org"),
			Contact: u.R("contact@example.org"),
		},
		Phone: u.R("+47-12345678"),
		Address: &AddressesDto{
			Visit: &VisitAddress{
				Street:             u.R("Torgallmenningen 1"),
				PostalArea:         u.R("Bergen"),
				PostalCode:         u.R("5009"),
				CountryCode:        u.R("NO"),
				CountrySubdivision: u.R("Vestland"),
			},
		},
		IsSupplier:         u.R(true),
		OrganizationNumber: u.R("123456789"),
	}
	err = cPostRequest.FromCustomerPostRequest0(CustomerPostRequest0{
		IsCompany: CustomerPostRequest0IsCompany(true),
		Name:      "ACME INC.",
	})
	require.NoError(err)
	resCustomer, err := c.CreateCustomerWithResponse(t.Context(), cPostRequest)
	require.NoError(err)
	require.NotNil(resCustomer.JSON200, "no customer was created")

	orderPostRequest := PostSalesordersJSONRequestBody{
		Status:       u.R(SalesOrderStatusEnumDraft),
		InternalMemo: u.R("some internal memo"),
		Memo:         u.R("a customer facing memo"),
		Customer: &struct {
			City                  *string          "json:\"city,omitempty\""
			CountryCode           *string          "json:\"countryCode,omitempty\""
			CountrySubdivision    *string          "json:\"countrySubdivision,omitempty\""
			Gln                   *string          "json:\"gln,omitempty\""
			Id                    int              "json:\"id\""
			InvoiceEmailAddresses *[]types.Email   "json:\"invoiceEmailAddresses,omitempty\""
			Name                  string           "json:\"name\""
			OrganizationNumber    *string          "json:\"organizationNumber,omitempty\""
			PostalArea            *string          "json:\"postalArea,omitempty\""
			PostalCode            *string          "json:\"postalCode,omitempty\""
			Street                *MultilineString "json:\"street,omitempty\""
		}{
			Id:   int(*resCustomer.JSON200.Id),
			Name: "CoolCustomer: what we know of them at the time of purchase",
		},
	}
	orderRes, err := c.PostSalesordersWithResponse(t.Context(), orderPostRequest)
	require.NoError(err)
	require.NotNil(orderRes.JSON200, "no order created")

	salesLineRes, err := c.PostSalesordersIdLinesWithResponse(t.Context(), int32(u.D(orderRes.JSON200.Id)), LineWithoutId{ //nolint:gosec
		Type: u.R(LineTypeEnumProduct),
		Product: &Product{
			Id: resProduct.JSON201.Id,
		},
		Description: u.R("Badeball!"),
		Price:       u.R(float32(100)),
		Quantity:    u.R(float32(10)),
	})
	require.NoError(err)
	require.NotNil(salesLineRes.JSON200, "no salesLine added")
}
