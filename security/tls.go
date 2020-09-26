package security

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

func NewTLSConfigWithClientCertificates() (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certPool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("security/mosquitto.org.crt")
	if err != nil {
		return nil, err
	}

	if err == nil {
		certPool.AppendCertsFromPEM(pemCerts)
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair("security/client-crt.pem", "security/client-key.pem")
	if err != nil {
		return nil, err
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certPool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.RequestClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}

func NewTLSConfig() (*tls.Config, error) {
	// Get the SystemCertPool, continue with an empty pool on error
	roots, _ := x509.SystemCertPool()
	if roots == nil {
		roots = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile("security/mosquitto.org.crt")
	if err != nil {
		return nil, err
	}

	if ok := roots.AppendCertsFromPEM(certs); !ok {
		return nil, errors.New("failed to append to root CA")
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS10,
		RootCAs:            roots,
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: false,
	}, nil
}
