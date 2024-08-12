// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package option

import (
	"net/http"
	"os"

	"go.mondoo.com/mondoo-go/internal"
	"go.mondoo.com/mondoo-go/internal/signer"
	"golang.org/x/oauth2"
)

const (
	endpointUS = "https://us.api.mondoo.com/query"
	endpointEU = "https://eu.api.mondoo.com/query"
)

// ClientOption is a configuration option for a Mondoo GraphQL API client.
type ClientOption interface {
	Apply(*internal.DialSettings)
}

// WithUserAgent returns a ClientOption that specifies the user agent to use.
// It is incompatible with WithHTTPClient.
func WithUserAgent(ua string) ClientOption {
	return withUA(ua)
}

type withUA string

func (w withUA) Apply(o *internal.DialSettings) { o.UserAgent = string(w) }

// WithHTTPClient returns a ClientOption that specifies the http.Client to use
func WithHTTPClient(client *http.Client) ClientOption {
	return withHTTPClient{client}
}

type withHTTPClient struct{ client *http.Client }

func (w withHTTPClient) Apply(o *internal.DialSettings) {
	o.HTTPClient = w.client
}

// WithTokenSource returns a ClientOption that specifies the oauth2.TokenSource
func WithTokenSource(s oauth2.TokenSource) ClientOption {
	return withTokenSource{"", s, nil}
}

type withTokenSource struct {
	endpoint string
	ts       oauth2.TokenSource
	err      error
}

func (w withTokenSource) Apply(o *internal.DialSettings) {
	o.TokenSource = w.ts
	o.TokenError = w.err
	if w.endpoint != "" {
		o.Endpoint = w.endpoint + "/query"
	}
}

// WithAPIToken returns a ClientOption that specifies the oauth2.TokenSource with the given token.
func WithAPIToken(token string) ClientOption {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return withTokenSource{"", src, nil}
}

// WithServiceAccount returns a ClientOption that specifies the credentials to use.
func WithServiceAccount(data []byte) ClientOption {
	ts, sa, err := signer.NewServiceAccountTokenSource(data)
	if err != nil {
		return withTokenSource{"", nil, err}
	}
	return withTokenSource{sa.ApiEndpoint, ts, nil}
}

// WithServiceAccountFile returns a ClientOption that specifies the credentials file to use.
func WithServiceAccountFile(filename string) ClientOption {
	data, err := os.ReadFile(filename)
	if err != nil {
		return withTokenSource{"", nil, err}
	}
	return WithServiceAccount(data)
}

// WithoutAuthentication returns a ClientOption that disables authentication.
// It is incompatible with WithTokenSource.
func WithoutAuthentication() ClientOption {
	return withoutAuthentication{}
}

type withoutAuthentication struct{}

func (w withoutAuthentication) Apply(o *internal.DialSettings) { o.NoAuth = true }

// WithEndpoint returns a ClientOption that specifies the endpoint to use.
func WithEndpoint(url string) ClientOption {
	return withEndpoint(url)
}

type withEndpoint string

func (w withEndpoint) Apply(o *internal.DialSettings) {
	o.Endpoint = string(w)
}

// UseUSRegion returns a ClientOption that specifies the US region endpoint.
func UseUSRegion() ClientOption {
	return withEndpoint(endpointUS)
}

// UseEURegion returns a ClientOption that specifies the EU region endpoint.
func UseEURegion() ClientOption {
	return withEndpoint(endpointEU)
}

// WithDefaultEndpoint returns a ClientOption that specifies the default region endpoint.
func WithDefaultEndpoint() ClientOption {
	return UseUSRegion()
}
