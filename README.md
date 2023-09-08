# Mondoo Platform GraphQL API Client for Go

Status: It is currently in beta.

This is a Go client library for the Mondoo Platform API using GraphQL. Mondoo is a security, risk and compliance
platform that provides continuous insight into your infrastructure. This library enables you to interact
programmatically with the Mondoo Platform to perform various tasks such as querying assets, setup integrations, and
fetching vulnerability and policy reports.

## Features

- Easy-to-use API to query data from Mondoo
- Strongly typed GraphQL queries and mutations
- Support for GraphQL subscriptions
- Token-based authentication

## Requirements

Our libraries are compatible with at least the three most recent, major Go
releases. They are currently compatible with:

- Go 1.21
- Go 1.20
- Go 1.19

In addition, a Mondoo account with API access is required to use this library.

## Installation

To install the package, run:

```bash
go get go.mondoo.com/mondoo-go
```

## Quick Start

Here is a quick example to fetch the details of an asset:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mondoo.com/mondoo-go"
	"go.mondoo.com/mondoo-go/option"
)

func main() {
	// Initialize the client
	client, err := mondoogql.NewClient(option.UseUSRegion(), option.WithAPIToken(os.Getenv("MONDOO_API_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch asset details
	assetMrn := "//assets.api.mondoo.app/spaces/dreamy-ellis-859675/assets/2TwUNCJcoPG5vHfUJaMf2gRgIaY"
	var q struct {
		Report struct {
			AssetReport struct {
				Asset struct {
					Name string
				}
			} `graphql:"... on AssetReport"`
		} `graphql:"assetReport(input: { assetMrn: $assetMrn} )"`
	}
	variables := map[string]interface{}{
		"assetMrn": mondoogql.ID(assetMrn),
	}
	err = client.Query(context.Background(), &q, variables)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(q.Report.AssetReport.Asset.Name)
}
```

## Test

Run all tests:

```bash
make test
```

For any requests, bug or comments, please [open an issue][issues] or [submit a
pull request][pulls].

## Kudos

This implementation is heavily inspired by the [GitHub GraphQL Go Client](https://github.com/shurcooL/githubv4).

[issues]: https://github.com/mondoohq/mondoo-go/issues/new

[pulls]: https://github.com/mondoohq/mondoo-go/pulls