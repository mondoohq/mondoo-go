// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package signer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	subject "go.mondoo.com/mondoo-go/internal/signer"
)

func TestNewServiceAccountTokenSource(t *testing.T) {
	t.Run("Invalid Data", func(t *testing.T) {
		data := []byte("invalid-yaml-data")

		tokenSource, creds, err := subject.NewServiceAccountTokenSource(data)

		assert.Nil(t, tokenSource)
		assert.Nil(t, creds)
		assert.Error(t, err)
		assert.Equal(t, "valid service account needs to be provided", err.Error())
	})

	t.Run("Invalid Private Key", func(t *testing.T) {
		credentials := []byte(`
certificate: |
    -----BEGIN CERTIFICATE-----
		foo
    -----END CERTIFICATE-----
force: false
mrn: //test.api.mondoo.app/spaces/test-796596/serviceaccounts/abc
private_key: |
    invalid-pem-key
space_mrn: //captain.api.mondoo.app/spaces/test-796596
`)

		tokenSource, creds, err := subject.NewServiceAccountTokenSource(credentials)

		assert.Nil(t, tokenSource)
		assert.Nil(t, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valid service account needs to be provided")
	})

	t.Run("Missing Private Key in Credentials in YAML format", func(t *testing.T) {
		credentials := []byte(`
mrn: //test.api.mondoo.app/spaces/test-796596/serviceaccounts/abc
space_mrn: //captain.api.mondoo.app/spaces/test-796596
`)

		tokenSource, creds, err := subject.NewServiceAccountTokenSource(credentials)
		assert.Nil(t, tokenSource)
		assert.Nil(t, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot load retrieved key")
	})
}
