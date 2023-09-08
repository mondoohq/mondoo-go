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
