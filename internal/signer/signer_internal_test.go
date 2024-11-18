// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package signer

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyFromBytes(t *testing.T) {
	t.Run("Valid ECDSA Private Key", func(t *testing.T) {
		privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		privKeyBytes, _ := x509.MarshalPKCS8PrivateKey(privKey)
		pemBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privKeyBytes})

		key, err := privateKeyFromBytes(pemBlock)
		assert.NoError(t, err)
		assert.NotNil(t, key)
		assert.IsType(t, &ecdsa.PrivateKey{}, key)
	})

	t.Run("Invalid PEM Format", func(t *testing.T) {
		_, err := privateKeyFromBytes([]byte("invalid-pem"))
		assert.ErrorIs(t, err, ErrAuthKeyNotPem)
	})

	t.Run("Invalid Private Key Type", func(t *testing.T) {
		// Generate an RSA private key (unsupported for this function)
		rsaKey, _ := x509.MarshalPKCS8PrivateKey(&struct{}{})
		pemBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: rsaKey})

		_, err := privateKeyFromBytes(pemBlock)
		assert.ErrorContains(t, err, "syntax error: sequence truncated")
	})
}
