// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// DialSettings holds information required to establish a connection
// with a Mondoo GraphQL API.
type DialSettings struct {
	SkipValidation bool
	NoAuth         bool
	Endpoint       string
	TokenSource    oauth2.TokenSource
	UserAgent      string
	HTTPClient     *http.Client
}

// Validate checks if dial settings are invalid.
func (ds *DialSettings) Validate() error {
	if ds.SkipValidation {
		return nil
	}
	hasCreds := ds.TokenSource != nil
	if ds.NoAuth && hasCreds {
		return errors.New("cannot use both WithoutAuthentication and WithTokenSource in combination")
	}
	return nil
}
