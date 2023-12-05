// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build debugtest
// +build debugtest

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	token := os.Getenv("MONDOO_API_TOKEN")
	err := generateSchema(token, "..")
	require.NoError(t, err)
}
func TestGetAPIEndpoint(t *testing.T) {
	// Test case 1: MONDOO_API_ENDPOINT not set
	expectedEndpoint := "https://us.api.mondoo.com/query"
	expectedHost := "us.api.mondoo.com"
	actualEndpoint, actualHost := getAPIEndpoint()
	require.Equal(t, expectedEndpoint, actualEndpoint, "Unexpected API endpoint")
	require.Equal(t, expectedHost, actualHost, "Unexpected host")

	// Test case 2: MONDOO_API_ENDPOINT set
	os.Setenv("MONDOO_API_ENDPOINT", "custom.api.endpoint")
	expectedEndpoint = "https://custom.api.endpoint/query"
	expectedHost = "custom.api.endpoint"
	actualEndpoint, actualHost = getAPIEndpoint()
	require.Equal(t, expectedEndpoint, actualEndpoint, "Unexpected API endpoint")
	require.Equal(t, expectedHost, actualHost, "Unexpected host")

	os.Setenv("MONDOO_API_ENDPOINT", "http://custom.api.endpoint:1234")
	expectedEndpoint = "http://custom.api.endpoint:1234/query"
	expectedHost = "custom.api.endpoint:1234"
	actualEndpoint, actualHost = getAPIEndpoint()
	require.Equal(t, expectedEndpoint, actualEndpoint, "Unexpected API endpoint")
	require.Equal(t, expectedHost, actualHost, "Unexpected host")
}
