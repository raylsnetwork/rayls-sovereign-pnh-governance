// Package mtls builds *tls.Config values for mutual-TLS connections
// to rayls services (currently NATS). The package is a portable copy
// of the equivalent helper in the rayls-privacy-relayer-api repo —
// duplicated rather than imported across repos.
package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// LoadServerConfig builds a *tls.Config for a server accepting mTLS
// connections. It loads the server keypair, trusts the supplied CA for
// verifying incoming client certs, and requires every client to present
// a cert signed by that CA.
func LoadServerConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading server keypair (%s, %s): %w", certFile, keyFile, err)
	}

	pool, err := loadCAPool(caFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// LoadClientConfig builds a *tls.Config for a peer dialling another
// rayls service (NATS, gRPC, etc.). It loads the client keypair,
// trusts the supplied CA for verifying the server certificate, and
// lets the dial target's hostname drive SAN verification.
func LoadClientConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading client keypair (%s, %s): %w", certFile, keyFile, err)
	}

	pool, err := loadCAPool(caFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

func loadCAPool(caFile string) (*x509.CertPool, error) {
	pem, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("reading CA file %s: %w", caFile, err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("CA file %s contains no PEM-encoded certificates", caFile)
	}
	return pool, nil
}
