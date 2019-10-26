package signingserver

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
)

var caCertPem []byte
var caKeyPem []byte
var caSerial *big.Int

func CAInit() {
	caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	notBefore := time.Now()
	notAfter := notBefore.Add(24 * 365 * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	caSerial = serialNumber
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Root CA"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, caKey.Public(), caKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	certBuf := bytes.Buffer{}
	if err = pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert: %s", err)
	}
	caCertPem = certBuf.Bytes()

	privBytes, err := x509.MarshalPKCS8PrivateKey(caKey)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	keyBuf := bytes.Buffer{}
	if err := pem.Encode(&keyBuf, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key: %s", err)
	}

	caKeyPem = keyBuf.Bytes()
}

func signCSRWithLocalCA(csr []byte) ([]byte, error) {
	block, _ := pem.Decode(csr)
	csrTemplate, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}
	serial := new(big.Int)
	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour)
	template := x509.Certificate{
		SerialNumber:          serial.Add(caSerial, big.NewInt(1)),
		Subject:               csrTemplate.Subject,
		SignatureAlgorithm:    csrTemplate.SignatureAlgorithm,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	pemBlockPub, _ := pem.Decode(caCertPem)
	caTemplate, err := x509.ParseCertificate(pemBlockPub.Bytes)
	if err != nil {
		return nil, err
	}
	pemBlockPriv, _ := pem.Decode(caKeyPem)
	caKey, err := x509.ParsePKCS8PrivateKey(pemBlockPriv.Bytes)
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caTemplate, csrTemplate.PublicKey, caKey)
	if err != nil {
		return nil, err
	}
	certBuf := bytes.Buffer{}
	if err = pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert: %s", err)
	}
	return certBuf.Bytes(), nil
}

