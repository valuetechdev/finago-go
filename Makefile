.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

.PHONY: generate
generate:
	go generate ./...
	go fmt ./...
	sed -i '' 's/int \[\]/Int \[\]/g' ./soapy/**/*.gen.go
	sed -i '' 's/IsCompany = "true"/IsCompany = true/g' ./resty/resty.gen.go
	sed -i '' 's/IsCompany = "false"/IsCompany = false/g' ./resty/resty.gen.go

.PHONY: api-payday
api-payday:
	@echo "openapi: getting latest Swagger from Payday and converting to OpenAPI 3"
	@curl https://converter.swagger.io/api/convert?url=https://me.24sevenoffice.com/swagger.json | jq > ./api/openapi/payroll.json

.PHONY: api-rest
api-rest:
	@echo "openapi: getting latest Finago Office's REST API"
	@curl https://rest-api.developer.24sevenoffice.com/doc/v1.yaml > ./api/openapi/rest.yaml

.PHONY: api
api:
	make -j2 api-payday api-rest	
	
.PHONY: check
check:
	go vet ./...
	go tool -modfile=go.tool.mod golangci-lint run ./...

# Bumps patch version
.PHONY: bump
bump:
	go tool -modfile=go.tool.mod git-bump 
