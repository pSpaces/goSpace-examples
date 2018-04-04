package certificate

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"
)

func GenerateCertConfigs() (*tls.Config, *tls.Config) {
	rootCert, err := generateRootCert("temporary root", time.Now().AddDate(1, 0, 0))
	if err != nil {
		log.Fatal(err)
	}

	// Generate client certificate.
	clientCert, err := generateCert(rootCert)
	if err != nil {
		log.Fatal(err)
	}
	if err := verify(clientCert.Cert, rootCert.Cert); err != nil {
		log.Fatalf("initial verification clientcert failed: %v", err)
	}

	// Generate server certificate.
	serverCert, err := generateCert(rootCert)
	if err != nil {
		log.Fatal(err)
	}
	if err := verify(serverCert.Cert, rootCert.Cert); err != nil {
		log.Fatalf("initial verification of servercert failed: %v", err)
	}

	serverConfig := server(serverCert.CertPEM(), serverCert.PrivateKeyPEM(), rootCert)
	if err != nil {
		log.Printf("server error: %v", err)
	}

	clientConfig := client(rootCert.CertPEM(), clientCert.CertPEM(), clientCert.PrivateKeyPEM())
	if err != nil {
		log.Printf("client error: %v", err)
		select {}
	}

	return serverConfig, clientConfig
}

func verify(clientCert, rootCert *x509.Certificate) error {
	pool := x509.NewCertPool()
	pool.AddCert(rootCert)
	opts := x509.VerifyOptions{
		DNSName: "localhost",
		Roots:   pool,
	}
	_, err := clientCert.Verify(opts)
	return err
}

// client dials the given address using the appropriate
// security primitives.
func client(rootCertPEM, certPEM, privPEM []byte) *tls.Config {
	pool, _ := certPool(rootCertPEM)
	cert, _ := tls.X509KeyPair(certPEM, privPEM)
	// if err != nil {
	// 	return fmt.Errorf("cannot make cert pool: %v", err)
	// }
	clientConfig := &tls.Config{
		Rand:               randReader,
		Certificates:       []tls.Certificate{cert},
		ServerName:         "localhost",
		RootCAs:            pool,
		InsecureSkipVerify: false,
	}

	return clientConfig

}

// certPool returns a pool of all the certificates in
// the given PEM data.
func certPool(certPEM []byte) (*x509.CertPool, error) {
	certs, err := ParsePEMCertificates(certPEM)
	if err != nil {
		return nil, err
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found in cert data")
	}
	//fmt.Printf("read certificates:\n%s\n", jsonString(certs))
	pool := x509.NewCertPool()
	for _, cert := range certs {
		pool.AddCert(cert)
	}
	return pool, nil
}

// server represents the server side of the protocol.
// The server accepts TLS connections and starts
// an echo server for each connection.
// The two strings, representing files containing
// PEM-format data, hold the private key and
// certificate of the server.
//func server(certPEM, privPEM []byte) (addr string, err error) {
func server(certPEM, privPEM []byte, rootCert *CertAndKey) *tls.Config {
	pool, errpool := certPool(rootCert.CertPEM())
	if errpool != nil {
		log.Fatalf("cannot create cert pool: %v", errpool)
	}
	cert, err := tls.X509KeyPair(certPEM, privPEM)
	if err != nil {
		log.Fatalf("cannot parse key pair: %v", err)
	}

	cfg := &tls.Config{
		Rand:               randReader,
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          pool,
		ServerName:         "localhost",
		RootCAs:            pool,
		InsecureSkipVerify: false,
	}

	return cfg
}

func createCertificate(template, parent *x509.Certificate, pub, priv interface{}) (cert *x509.Certificate, derBytes []byte, err error) {
	derBytes, err = x509.CreateCertificate(randReader, template, parent, pub, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("canot create certificate: %v", err)
	}
	certs, err := x509.ParseCertificates(derBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse new certificate: %v", err)
	}
	if len(certs) != 1 {
		return nil, nil, fmt.Errorf("need exactly one certificate")
	}
	return certs[0], derBytes, nil
}

// ParsePEMCertificates parses PEM/DER encoded certificates from
// the given PEM data.
func ParsePEMCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		var der *pem.Block
		der, pemData = pem.Decode(pemData)
		if der == nil {
			break
		}
		if der.Type == "CERTIFICATE" {
			dcerts, err := x509.ParseCertificates(der.Bytes)
			if err != nil {
				return nil, err
			}
			certs = append(certs, dcerts...)
		}
	}
	return certs, nil
}

// CertAndKey holds a certificate and its principal private
// key.
type CertAndKey struct {
	Cert        *x509.Certificate
	certDERData []byte
	PrivateKey  *rsa.PrivateKey
}

// CertPEM returns the certificate data in PEM format.
func (ck *CertAndKey) CertPEM() []byte {
	var b bytes.Buffer
	pem.Encode(&b, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ck.certDERData,
	})
	return b.Bytes()
}

// PrivateKeyPEM returns the private key data in PEM format.
func (ck *CertAndKey) PrivateKeyPEM() []byte {
	var b bytes.Buffer
	pem.Encode(&b, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ck.PrivateKey),
	})
	return b.Bytes()
}

// generateRoot generates a self-signed root certificate
// and returns the cert/key pair.
func generateRootCert(name string, expiry time.Time) (*CertAndKey, error) {
	priv, err := rsa.GenerateKey(randReader, 512)
	if err != nil {
		return nil, fmt.Errorf("cannot generate key: %v", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			CommonName:   name,
			Organization: []string{"juju"},
		},
		NotBefore: now.Add(-5 * time.Minute).UTC(),
		NotAfter:  expiry,

		SubjectKeyId:          bigIntHash(priv.N),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:       true,
		MaxPathLen: 1,
	}
	cert, derBytes, err := createCertificate(&template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	return &CertAndKey{
		Cert:        cert,
		certDERData: derBytes,
		PrivateKey:  priv,
	}, nil
}

// generateCert generates a new private key, signs a certificate
// in the name "any" with the signer's private key and returns the
// cert/key pair.
func generateCert(signer *CertAndKey) (*CertAndKey, error) {
	priv, err := rsa.GenerateKey(randReader, 512)
	if err != nil {
		return nil, fmt.Errorf("cannot generate key: %v", err)
	}
	now := time.Now()
	template := x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"hmm"},
		},
		NotBefore: now.Add(-5 * time.Minute).UTC(),
		NotAfter:  now.AddDate(1, 0, 0).UTC(), // valid for 1 year.

		SubjectKeyId: bigIntHash(priv.N),
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature |
			x509.KeyUsageDataEncipherment,
	}
	cert, derBytes, err := createCertificate(&template, signer.Cert, &priv.PublicKey, signer.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &CertAndKey{
		Cert:        cert,
		certDERData: derBytes,
		PrivateKey:  priv,
	}, nil
}

func bigIntHash(n *big.Int) []byte {
	h := sha1.New()
	h.Write(n.Bytes())
	return h.Sum(nil)
}

// The randReader is necessary for the code to work.
type pseudoRand struct{}

//randReader is a replacement for rand.Reader.
var randReader = pseudoRand{}

func (pseudoRand) Read(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = byte(rand.Int())
	}
	return len(buf), nil
}
