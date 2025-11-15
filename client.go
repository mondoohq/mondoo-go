// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mondoogql

import (
	"context"
	"net/http"

	"github.com/shurcooL/graphql"
	"go.mondoo.com/mondoo-go/internal"
	internal_http "go.mondoo.com/mondoo-go/internal/http"
	"go.mondoo.com/mondoo-go/option"
)

// Client is a Mondoo GraphQL API client.
type Client struct {
	client *graphql.Client
}

// NewClient creates a new Mondoo GraphQL API client
func NewClient(opts ...option.ClientOption) (*Client, error) {
	// NOTE: we need to prepend options, so we don't override user-specified options
	opts = append([]option.ClientOption{option.WithDefaultEndpoint(), option.WithUserAgent("mondoo-graphql-client/" + internal.Version)}, opts...)
	httpClient, endpoint, err := internal_http.NewHttpClient(opts...)
	if err != nil {
		return nil, err
	}

	c := &Client{}
	c.client = graphql.NewClient(endpoint, httpClient)
	return c, nil
}

// Query executes a single GraphQL query request, populating the response into q. It supports
// two arguments for the query and variables.
// - q should be a pointer to struct that corresponds to the Mondoo GraphQL schema.
// - variables should be a map of string to arbitrary Go type that corresponds to the GraphQL variables.
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.client.Query(ctx, q, variables)
}

// Mutate executes a single GraphQL mutation request, populating the response into m. It supports
// three arguments for the mutation, input and variables.
// - m should be a pointer to struct that corresponds to the Mondoo GraphQL schema.
// - input should be a map of string to arbitrary Go type that corresponds to the GraphQL input.
// - variables should be a map of string to arbitrary Go type that corresponds to the GraphQL variables.
func (c *Client) Mutate(ctx context.Context, m interface{}, input Input, variables map[string]interface{}) error {
	if input != nil {
		if variables == nil {
			variables = map[string]interface{}{"input": input}
		} else {
			variables["input"] = input
		}
	}
	return c.client.Mutate(ctx, m, variables)
}

// NewHttpClient creates a new *http.Client instance based on the given options.
func NewHttpClient(opts ...option.ClientOption) (*http.Client, string, error) {
	return internal_http.NewHttpClient(opts...)
}
