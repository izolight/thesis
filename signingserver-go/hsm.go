package signingserver

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// temporary key storage -> very unsafe
var keyStorage map[string]*rsa.PrivateKey

func generateLocalKey(template *x509.CertificateRequest) ([]byte, error) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyStorage[template.Subject.CommonName] = keyBytes
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, keyBytes)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	err = pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	return buf.Bytes(), err
}

func sign(hash Hash, subject string) ([]byte, error){
	key, ok := keyStorage[subject]
	if !ok {
		return nil, fmt.Errorf("Key with subject %s not found", subject)
	}
	return key.Sign(rand.Reader, []byte(hash), crypto.SHA256)
}