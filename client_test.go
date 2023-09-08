// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mondoogql

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mondoo.com/mondoo-go/option"
)

func TestGraphQLClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/query" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		requestQuery, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		q := strings.TrimSpace(string(requestQuery))
		var response string
		if q == "{\"query\":\"{viewer{mrn}}\"}" {
			response = `{"data":{"viewer":{"mrn":"//api.mondoo.app/v1/serviceaccount/1234567890"}}}`
		}
		if response == "" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	client, err := NewClient(option.WithEndpoint(ts.URL + "/query"))
	require.NoError(t, err)

	var q struct {
		Viewer struct {
			Mrn string
		}
	}
	err = client.Query(context.Background(), &q, nil)
	require.NoError(t, err)
	assert.Equal(t, "//api.mondoo.app/v1/serviceaccount/1234567890", q.Viewer.Mrn)
}
