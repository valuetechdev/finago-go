//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hooklift/gowsdl"
)

type Client struct {
	packagename string
	wsdlPath    string
}

const WSDL_BASE_PATH = "../api/wsdl/"

var clientsToGenerate = []Client{
	{packagename: "account", wsdlPath: WSDL_BASE_PATH + "accountService.wsdl"},
	{packagename: "attachment", wsdlPath: WSDL_BASE_PATH + "attachmentService.wsdl"},
	{packagename: "auth", wsdlPath: WSDL_BASE_PATH + "authService.wsdl"},
	{packagename: "client", wsdlPath: WSDL_BASE_PATH + "clientService.wsdl"},
	{packagename: "company", wsdlPath: WSDL_BASE_PATH + "companyService.wsdl"},
	{packagename: "invoice", wsdlPath: WSDL_BASE_PATH + "invoiceService.wsdl"},
	{packagename: "person", wsdlPath: WSDL_BASE_PATH + "personService.wsdl"},
	{packagename: "product", wsdlPath: WSDL_BASE_PATH + "productService.wsdl"},
	{packagename: "project", wsdlPath: WSDL_BASE_PATH + "projectService.wsdl"},
	{packagename: "transaction", wsdlPath: WSDL_BASE_PATH + "transactionService.wsdl"},
}

const PKG_PATH = "."

func main() {
	for _, c := range clientsToGenerate {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		gw, err := gowsdl.NewGoWSDL(string(c.wsdlPath), c.packagename, false, true)
		if err != nil {
			fmt.Print(
				fmt.Errorf("could not generate package=%s", c.packagename),
			)
			return
		}

		res, _ := gw.Start()

		pkgAndFileName := fmt.Sprintf("%s/%s.gen.go", c.packagename, c.packagename)
		path, _ := filepath.Abs(fmt.Sprintf("%s/%s", cwd, pkgAndFileName))

		header, _ := res["header"]
		types, _ := res["types"]
		operations, _ := res["operations"]

		dir := fmt.Sprintf("%s/%s", cwd, c.packagename)
		if err = os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}

		if _, err := os.ReadFile(path); err != nil {
			os.Create(path)
		}

		err = os.WriteFile(path, []byte(string(header)+string(types)+string(operations)), 0646)
		if err != nil {
			fmt.Printf("%s\n", err)
			panic(err)
		}

	}
}
