// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"errors"
	"net/http"

	"go.mondoo.com/mondoo-go/internal"
	"go.mondoo.com/mondoo-go/option"
	"golang.org/x/oauth2"
)

const (
	userAgent = "User-Agent"
)

// newSettings applies the given options to a DialSettings struct.
func newSettings(opts []option.ClientOption) (*internal.DialSettings, error) {
	var o internal.DialSettings
	for _, opt := range opts {
		opt.Apply(&o)
	}
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return &o, nil
}

// NewHttpClient creates a new *http.Client instance based on the given options.
func NewHttpClient(opts ...option.ClientOption) (*http.Client, string, error) {
	settings, err := newSettings(opts)
	if err != nil {
		return nil, "", err
	}

	trans, err := newTransport(http.DefaultTransport, settings)
	if err != nil {
		return nil, "", err
	}

	return &http.Client{Transport: trans}, settings.Endpoint, nil
}

// newTransport creates a new http.RoundTripper based on the given options.
func newTransport(base http.RoundTripper, settings *internal.DialSettings) (http.RoundTripper, error) {
	paramTransport := &parameterTransport{
		base:      base,
		userAgent: settings.UserAgent,
	}

	// configure authentication
	var trans http.RoundTripper = paramTransport
	switch {
	case settings.NoAuth:
		// Do nothing.
	case settings.TokenSource != nil:
		trans = &oauth2.Transport{
			Base:   trans,
			Source: settings.TokenSource,
		}
	default:
		return nil, errors.New("no authentication configured")
	}
	return trans, nil
}

// parameterTransport is an http.RoundTripper that adds parameters to outgoing requests.
type parameterTransport struct {
	userAgent string
	base      http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t *parameterTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.base
	if rt == nil {
		return nil, errors.New("transport: no Transport specified")
	}
	newReq := *req
	newReq.Header = make(http.Header)
	for k, vv := range req.Header {
		newReq.Header[k] = vv
	}
	if t.userAgent != "" {
		newReq.Header.Set(userAgent, t.userAgent)
	}

	return rt.RoundTrip(&newReq)
}
