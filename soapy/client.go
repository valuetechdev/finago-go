//go:generate go run generate_services.go

package soapy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hooklift/gowsdl/soap"
	"github.com/valuetechdev/finago-go/soapy/account"
	"github.com/valuetechdev/finago-go/soapy/attachment"
	"github.com/valuetechdev/finago-go/soapy/auth"
	clientservice "github.com/valuetechdev/finago-go/soapy/client"
	"github.com/valuetechdev/finago-go/soapy/company"
	"github.com/valuetechdev/finago-go/soapy/invoice"
	"github.com/valuetechdev/finago-go/soapy/person"
	"github.com/valuetechdev/finago-go/soapy/product"
	"github.com/valuetechdev/finago-go/soapy/project"
	"github.com/valuetechdev/finago-go/soapy/transaction"
)

const (
	accountURL     = "https://api.24sevenoffice.com/Economy/Account/V004/Accountservice.asmx"
	attachmentURL  = "https://webservices.24sevenoffice.com/Economy/Accounting/Accounting_V001/AttachmentService.asmx"
	authURL        = "https://api.24sevenoffice.com/authenticate/v001/authenticate.asmx"
	clientURL      = "https://api.24sevenoffice.com/Client/V001/ClientService.asmx"
	companyURL     = "https://api.24sevenoffice.com/CRM/Company/V001/CompanyService.asmx"
	invoiceURL     = "https://api.24sevenoffice.com/Economy/InvoiceOrder/V001/InvoiceService.asmx"
	personURL      = "https://webservices.24sevenoffice.com/CRM/Contact/PersonService.asmx"
	productURL     = "https://api.24sevenoffice.com/Logistics/Product/V001/ProductService.asmx"
	projectURl     = "https://webservices.24sevenoffice.com/Project/V001/ProjectService.asmx"
	transactionURL = "https://api.24sevenoffice.com/Economy/Accounting/V001/TransactionService.asmx"
)

type Client struct {
	sessionId   string
	headers     map[string]string
	credentials *auth.Credential
	httpClient  *http.Client
	didAuth     bool

	Account     account.AccountServiceSoap
	Attachment  attachment.AttachmentServiceSoap
	Auth        auth.AuthenticateSoap
	Client      clientservice.ClientServiceSoap
	Company     company.CompanyServiceSoap
	Invoice     invoice.InvoiceServiceSoap
	Person      person.PersonServiceSoap
	Product     product.ProductServiceSoap
	Project     project.ProjectServiceSoap
	Transaction transaction.TransactionServiceSoap
}

type Option func(*Client)

// WithHttpClient sets a custom http.Client. Defaults to [http.DefaultClient].
func WithHttpClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// authInterceptor wraps an http.RoundTripper to automatically check auth before requests
type authInterceptor struct {
	client    *Client
	transport http.RoundTripper
}

func (a *authInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// Skip auth check for auth service URLs to avoid recursion
	if !strings.EqualFold(req.URL.String(), authURL) {
		a.client.didAuth = true
		if err := a.client.CheckAuth(); err != nil {
			return nil, err
		}
	} else if !a.client.didAuth {
		a.client.didAuth = true
		if err := a.client.CheckAuth(); err != nil {
			return nil, err
		}
	}

	// Add current headers to the request
	for key, value := range a.client.headers {
		req.Header.Set(key, value)
	}

	return a.transport.RoundTrip(req)
}

// Returns new [Client].
//
// You can reuse an already generated sessionId by using [Client.SetSessionId]
func New(credentials auth.Credential, options ...Option) *Client {
	headers := map[string]string{}
	client := &Client{
		credentials: &credentials,
		headers:     headers,
		httpClient:  http.DefaultClient,
	}
	for _, option := range options {
		option(client)
	}

	// Wrap the httpClient transport with auth interceptor
	interceptedClient := &http.Client{
		Transport: &authInterceptor{
			client:    client,
			transport: client.httpClient.Transport,
		},
		Timeout: client.httpClient.Timeout,
	}
	if interceptedClient.Transport.(*authInterceptor).transport == nil {
		interceptedClient.Transport.(*authInterceptor).transport = http.DefaultTransport
	}
	client.httpClient = interceptedClient

	authService := auth.NewAuthenticateSoap(newSoapClient(authURL, client))
	invoiceService := invoice.NewInvoiceServiceSoap(newSoapClient(invoiceURL, client))
	productService := product.NewProductServiceSoap(newSoapClient(productURL, client))
	accountService := account.NewAccountServiceSoap(newSoapClient(accountURL, client))
	companyService := company.NewCompanyServiceSoap(newSoapClient(companyURL, client))
	clientService := clientservice.NewClientServiceSoap(newSoapClient(clientURL, client))
	personService := person.NewPersonServiceSoap(newSoapClient(personURL, client))
	projectService := project.NewProjectServiceSoap(newSoapClient(projectURl, client))
	transactionService := transaction.NewTransactionServiceSoap(newSoapClient(transactionURL, client))
	attachmentService := attachment.NewAttachmentServiceSoap(newSoapClient(attachmentURL, client))

	client.Account = accountService
	client.Attachment = attachmentService
	client.Auth = authService
	client.Client = clientService
	client.Company = companyService
	client.Invoice = invoiceService
	client.Person = personService
	client.Product = productService
	client.Project = projectService
	client.Transaction = transactionService

	return client
}

func newSoapClient(url string, client *Client) *soap.Client {
	return soap.NewClient(url, soap.WithHTTPClient(client.httpClient))
}

// Returns sessionId
func (c *Client) GetSessionId() string {
	return c.sessionId
}

// Sets sessionId
//
// And sets headers with "cookie: ASP.NET_SessionId=<sessionId>".
func (c *Client) SetSessionId(sessionId string) {
	c.headers["cookie"] = fmt.Sprintf("ASP.NET_SessionId=%s", sessionId)
	c.sessionId = sessionId
}

// Returns true if sessionId is valid
func (c *Client) IsSessionIdValid() bool {
	r, err := c.Auth.HasSession(&auth.HasSession{})
	if err != nil {
		return false
	}

	return r.HasSessionResult
}

// Checks if auth is valid.
//
// Reauthenticates if auth is invalid.
func (c *Client) CheckAuth() error {
	if c.IsSessionIdValid() {
		return nil
	}
	res, err := c.Auth.Login(&auth.Login{
		Credential: c.credentials,
	})
	if err != nil {
		return fmt.Errorf("soap24: failed to login: %w", err)
	}

	sessionId := res.LoginResult
	c.SetSessionId(sessionId)

	return nil
}
