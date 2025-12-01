package soapy

import (
	"os"
	"testing"
	"time"

	"github.com/hooklift/gowsdl/soap"
	"github.com/stretchr/testify/require"
	"github.com/valuetechdev/finago-go/soapy/account"
	"github.com/valuetechdev/finago-go/soapy/auth"
	"github.com/valuetechdev/finago-go/soapy/client"
	"github.com/valuetechdev/finago-go/soapy/company"
	"github.com/valuetechdev/finago-go/soapy/invoice"
	"github.com/valuetechdev/finago-go/soapy/person"
	"github.com/valuetechdev/finago-go/soapy/product"
	"github.com/valuetechdev/finago-go/soapy/project"
)

var applicationId = auth.Guid(os.Getenv("TFSO_SOAP_APPLICATIONID"))
var credentials = auth.Credential{
	ApplicationId: &applicationId,
	Username:      os.Getenv("TFSO_SOAP_USERNAME"),
	Password:      os.Getenv("TFSO_SOAP_PASSWORD"),
}

func TestClientInitialization(t *testing.T) {
	require := require.New(t)

	c := New(credentials)
	require.True(c.IsSessionIdValid(), "client should have a valid sessionId checkAuth")
}

func TestClientInitializationWithAuthCalledFirst(t *testing.T) {
	require := require.New(t)

	c := New(credentials)
	_, err := c.Auth.GetIdentities(&auth.GetIdentities{})
	require.NoError(err)
	require.True(c.IsSessionIdValid(), "client should have a valid sessionId after first call")
	require.NoError(c.CheckAuth(), "client should be able to check auth after first call")
}

func TestServices(t *testing.T) {
	require := require.New(t)

	c := New(credentials)

	changedAfter := soap.CreateXsdDateTime(time.Now(), true)
	_, err := c.Account.GetAccountList(&account.GetAccountList{})
	require.NoError(err, "GetAccountList")
	_, err = c.Auth.GetIdentities(&auth.GetIdentities{})
	require.NoError(err, "GetIdentities")
	_, err = c.Client.GetClientInformation(&client.GetClientInformation{})
	require.NoError(err, "GetClientInformation")
	_, err = c.Company.GetCompanies(&company.GetCompanies{SearchParams: &company.CompanySearchParameters{
		ChangedAfter: changedAfter,
	}})
	require.NoError(err, "GetCompanies")
	_, err = c.Invoice.GetInvoices(&invoice.GetInvoices{SearchParams: &invoice.InvoiceSearchParameters{
		ChangedAfter: changedAfter,
	}})
	require.NoError(err, "GetInvoices")
	_, err = c.Person.GetPersons(&person.GetPersons{
		PersonSearch: &person.PersonSearchParameters{
			ChangedAfter: changedAfter,
		},
	})
	require.NoError(err, "GetPersons")

	state := person.TriStateNone
	_, err = c.Person.GetPersons(&person.GetPersons{
		PersonSearch: &person.PersonSearchParameters{
			ChangedAfter: changedAfter,
			Email:        "lol",
			IsEmployee:   &state,
		},
	})
	require.NoError(err, "GetPersons")
	_, err = c.Product.GetProducts(&product.GetProducts{
		SearchParams: &product.ProductSearchParameters{DateChanged: changedAfter},
	})
	require.NoError(err, "GetProducts")
	_, err = c.Project.GetProjectList(&project.GetProjectList{Ps: &project.ProjectSearch{ChangedAfter: changedAfter}})
	require.NoError(err, "GetProjectList")
}
