package verifier

import (
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mozilla.org/pkcs7"
)

type SignatureContainerVerifier struct {
	container   []byte
	data        chan SignatureData
	signerEmail chan string
	additionalCerts []*x509.Certificate
}

func NewSignatureContainerVerifier(c []byte, additionalCerts []*x509.Certificate) *SignatureContainerVerifier {
	return &SignatureContainerVerifier{
		container:   c,
		data:        make(chan SignatureData, 1),
		signerEmail: make(chan string, 1),
		additionalCerts:additionalCerts,
	}
}

func (s *SignatureContainerVerifier) Verify(verifyLTV bool) error {
	p7, err := pkcs7.Parse(s.container)
	if err != nil {
		return fmt.Errorf("could not decode signature container: %w", err)
	}
	signatureData := SignatureData{}
	if err := proto.Unmarshal(p7.Content, &signatureData); err != nil {
		return fmt.Errorf("could not unmarshal signature data: %w", err)
	}
	s.data <- signatureData
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("could not get system cert pool: %w", err)
	}
	for _, cert := range s.additionalCerts{
		certPool.AddCert(cert)
	}

	if err := p7.VerifyWithChain(certPool); err != nil {
		return fmt.Errorf("could not verify pcks7: %w", err)
	}

	if verifyLTV {
		l := LTVVerifier{
			Certs: p7.Certificates,
			//LTVData: s.container,
		}
		if err := l.Verify(); err != nil {
			return fmt.Errorf("verifyLTV information for signature is not valid: %w", err)
		}
	}

	for _, c := range p7.Certificates {
		if !c.IsCA {
			s.signerEmail <- c.EmailAddresses[0]
			break
		}
	}

	return nil
}

func (s *SignatureContainerVerifier) SignatureData() SignatureData {
	return <-s.data
}

func (s *SignatureContainerVerifier) SignerEmail() string {
	return <-s.signerEmail
}
