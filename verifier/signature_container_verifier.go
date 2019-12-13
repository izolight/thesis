package verifier

import (
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"go.mozilla.org/pkcs7"
	"time"
)

type SignatureContainerVerifier struct {
	container       []byte
	data            chan SignatureData
	signer          chan *x509.Certificate
	signingTime     chan time.Time
	additionalCerts []*x509.Certificate
	cfg             *Config
}

func NewSignatureContainerVerifier(c []byte, additionalCerts []*x509.Certificate, cfg Config) *SignatureContainerVerifier {
	cfg.Logger = cfg.Logger.WithField("verifier", "signature container")
	return &SignatureContainerVerifier{
		container:       c,
		data:            make(chan SignatureData, 1),
		signer:          make(chan *x509.Certificate, 1),
		signingTime:     make(chan time.Time, 1),
		additionalCerts: additionalCerts,
		cfg:             &cfg,
	}
}

func (s *SignatureContainerVerifier) Verify(verifyLTV bool) error {
	s.cfg.Logger.Info("started verifying")
	p7, err := pkcs7.Parse(s.container)
	if err != nil {
		return fmt.Errorf("could not decode signature container: %w", err)
	}
	s.cfg.Logger.Info("parsed pkcs7 signature container")
	signatureData := SignatureData{}
	if err := proto.Unmarshal(p7.Content, &signatureData); err != nil {
		return fmt.Errorf("could not unmarshal signature data: %w", err)
	}
	s.cfg.Logger.WithFields(log.Fields{
		"signatureLevel":         signatureData.SignatureLevel,
		"salted_document_hashes": fmt.Sprintf("%x", signatureData.SaltedDocumentHash),
		"hash_algorithm":         signatureData.HashAlgorithm,
		"mac_key":                fmt.Sprintf("%x", signatureData.MacKey),
		"mac_algorithm":          signatureData.MacAlgorithm,
	}).Info("decoded signature data")
	s.data <- signatureData
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("could not get system cert pool: %w", err)
	}
	for _, cert := range s.additionalCerts {
		certPool.AddCert(cert)
	}

	if err := p7.VerifyWithChain(certPool); err != nil {
		return fmt.Errorf("could not verify pcks7: %w", err)
	}
	s.cfg.Logger.Info("verified pkcs7 certificate chain")

	if verifyLTV {
		l := LTVVerifier{
			Certs: p7.Certificates,
			//LTVData: s.container,
		}
		if err := l.Verify(); err != nil {
			return fmt.Errorf("verifyLTV information for signature is not valid: %w", err)
		}
	}
	signer := p7.GetOnlySigner()
	signingTime := <-s.signingTime
	if signer.NotBefore.After(signingTime) {
		return fmt.Errorf("certificate was issued at %s, which is after the signing time %s", signer.NotBefore, signingTime)
	}
	s.signer <- signer
	s.cfg.Logger.WithFields(log.Fields{
		"subject":    signer.Subject,
		"email":      signer.EmailAddresses[0],
		"issuer":     signer.Issuer,
		"not_before": signer.NotBefore,
		"not_after":  signer.NotAfter,
	}).Info("pkcs7 signerEmail infos")

	s.cfg.Logger.Info("finished verifying")
	return nil
}

func (s *SignatureContainerVerifier) SignatureData() SignatureData {
	return <-s.data
}

func (s *SignatureContainerVerifier) Signer() *x509.Certificate {
	return <-s.signer
}

func (s *SignatureContainerVerifier) SendSigningTime(signingTime time.Time) {
	s.signingTime <- signingTime
}