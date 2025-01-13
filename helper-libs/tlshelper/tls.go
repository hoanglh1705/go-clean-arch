package tlshelper

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"os"
)

func NewServerTLSConfigFromBase64(certValue,
	keyValue,
	rootCACertValue string,
	clientAuthType tls.ClientAuthType,
) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	certPem, err := base64.StdEncoding.DecodeString(certValue)
	if err != nil {
		return nil, err
	}
	keyPem, err := base64.StdEncoding.DecodeString(keyValue)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := base64.StdEncoding.DecodeString(rootCACertValue)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	if clientAuthType != tls.NoClientCert {
		// tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientAuth = clientAuthType
		tlsConfig.ClientCAs = caCertPool
	}

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}

func NewClientTLSConfigFromBase64(
	certValue,
	keyValue,
	caCertValue string,
	insecureSkipVerify bool,
) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	certPem, err := base64.StdEncoding.DecodeString(certValue)
	if err != nil {
		return nil, err
	}
	keyPem, err := base64.StdEncoding.DecodeString(keyValue)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	if caCertValue != "" {
		caCert, err := base64.StdEncoding.DecodeString(caCertValue)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	if insecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}

func NewServerTLSConfigFromFile(certFile, keyFile, caCertFile string,
	clientAuthType tls.ClientAuthType) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	if clientAuthType != tls.NoClientCert {
		// tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientAuth = clientAuthType
		tlsConfig.ClientCAs = caCertPool
	}

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}

func NewClientTLSConfigFromFile(certFile, keyFile, caCertFile string,
	insecureSkipVerify bool) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	if insecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}

func GetX509CertificateFromTlsCertificate(tlsCertificate *tls.Certificate) *x509.Certificate {
	return tlsCertificate.Leaf
}

func GetPrivateKeyFromTlsCertificate(tlsCertificate *tls.Certificate) crypto.PrivateKey {
	return tlsCertificate.PrivateKey
}

func GetPublicKeyFromTlsCertificate(tlsCertificate *tls.Certificate) crypto.PublicKey {
	if len(tlsCertificate.Certificate) == 0 {
		return nil
	}

	cert, err := x509.ParseCertificate(tlsCertificate.Certificate[0])
	if err != nil {
		return nil
	}

	if cert != nil {
		return cert.PublicKey
	}

	return nil
}
