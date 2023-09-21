.PHONY: license license/headers/check license/headers/apply

generate:
	echo "Ensure the MONDOO_API_TOKEN environment variable is set."
	echo "Generating code..."
	go generate ./...

license: license/headers/check

license/headers/check:
	copywrite headers --plan

license/headers/apply:
	copywrite headers

test: test/go test/lint

test/go:
	go test -cover $(shell go list ./...)

test/lint: test/lint/golangci-lint/run

prep/tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: test/lint/golangci-lint/run
test/lint/golangci-lint/run: prep/tools
	golangci-lint --version
	golangci-lint run