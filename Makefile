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

test:
	go test -cover $(shell go list ./... | grep -v '/providers/')