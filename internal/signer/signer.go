// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package signer

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"time"

	jose "github.com/go-jose/go-jose/v3"
	jwt "github.com/go-jose/go-jose/v3/jwt"
	"golang.org/x/oauth2"
)

const serviceAccountIssuer = "mondoo/ams"

// .p8 errors file.
var (
	ErrAuthKeyNotPem   = errors.New("AuthKey must be a valid .p8 PEM file")
	ErrAuthKeyNotECDSA = errors.New("AuthKey must be of type ecdsa.PrivateKey")
	ErrAuthKeyNil      = errors.New("AuthKey was nil")
)

type serviceAccountCredentials struct {
	Mrn         string `json:"mrn,omitempty"`
	ParentMrn   string `json:"parent_mrn,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	ApiEndpoint string `json:"api_endpoint,omitempty"`
}

// privateKeyFromBytes loads a .p8 certificate from an in memory byte array and
// returns an *ecdsa.PrivateKey.
func privateKeyFromBytes(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, ErrAuthKeyNotPem
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pk := key.(type) {
	case *ecdsa.PrivateKey:
		return pk, nil
	default:
		return nil, ErrAuthKeyNotECDSA
	}
}

func NewServiceAccountTokenSource(data []byte) (*serviceAccountTokenSource, *serviceAccountCredentials, error) {
	var credentials *serviceAccountCredentials
	err := json.Unmarshal(data, &credentials)
	if credentials == nil || err != nil {
		return nil, nil, errors.New("valid service account needs to be provided")
	}

	// verify that we can read the private key
	privateKey, err := privateKeyFromBytes([]byte(credentials.PrivateKey))
	if err != nil {
		return nil, nil, errors.New("cannot load retrieved key: " + err.Error())
	}

	// configure authentication plugin, since the server only accepts authenticated calls
	cfg := tokenSourceConfig{
		PrivateKey: privateKey,
		Issuer:     serviceAccountIssuer,
		Kid:        credentials.Mrn,
		Subject:    credentials.Mrn,
	}

	return &serviceAccountTokenSource{
		cfg: cfg,
	}, credentials, nil
}

type tokenSourceConfig struct {
	Subject    string
	Issuer     string
	Kid        string
	PrivateKey *ecdsa.PrivateKey
}

type serviceAccountTokenSource struct {
	cfg tokenSourceConfig
}

// sign generates a short-lived JWT bearer token.
func (cap *serviceAccountTokenSource) sign() (*jwt.Claims, string, error) {
	issuedAt := time.Now()

	cl := jwt.Claims{
		Subject:   cap.cfg.Subject,
		Issuer:    cap.cfg.Issuer,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		NotBefore: jwt.NewNumericDate(issuedAt),
	}

	// valid for 60 seconds
	cl.Expiry = jwt.NewNumericDate(issuedAt.Add(time.Duration(60) * time.Second))
	sig, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES384,
		Key:       cap.cfg.PrivateKey,
	}, (&jose.SignerOptions{}).WithHeader("kid", cap.cfg.Kid).WithType("JWT"))
	if err != nil {
		return nil, "", err
	}

	bearer, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return nil, "", err
	}

	return &cl, bearer, nil
}

// Token returns a token or an error.
func (cts *serviceAccountTokenSource) Token() (*oauth2.Token, error) {
	claims, bearer, err := cts.sign()
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: bearer,
		Expiry:      claims.Expiry.Time(),
	}, nil
}
