package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func loadCert(caCert string) (*x509.Certificate, error) {
	data, err := os.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode pem in %s", caCert)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func newTLSConfig(caCert string) (*tls.Config, error) {
	cert, err := loadCert(caCert)
	if err != nil {
		return nil, err
	}
	result := tls.Config{}
	result.RootCAs = x509.NewCertPool()
	result.RootCAs.AddCert(cert)
	return &result, nil
}
