package verifier

import (
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"go.mozilla.org/pkcs7"
	"golang.org/x/crypto/ocsp"
	"time"
)

type SignatureContainerVerifier struct {
	container       []byte
	data            chan SignatureData
	signingCertData chan signingCertData
	signingTime     chan time.Time
	additionalCerts []*x509.Certificate
	cfg             *Config
}

type signingCertData struct {
	Signer      string      `json:"signer"`
	SignerEmail string      `json:"signer_email"`
	Certs       []CertChain `json:"cert_chain"`
}

func NewSignatureContainerVerifier(c []byte, additionalCerts []*x509.Certificate, cfg Config) *SignatureContainerVerifier {
	cfg.Logger = cfg.Logger.WithField("verifier", "signature container")
	return &SignatureContainerVerifier{
		container:       c,
		data:            make(chan SignatureData, 1),
		signingTime:     make(chan time.Time, 1),
		signingCertData: make(chan signingCertData),
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

	ocspStatus := make(map[string]*ocsp.Response)
	if verifyLTV {
		l, err := NewLTVVerifier(p7.Certificates, p7.CRLs, p7.RawOCSPResponses)
		if err != nil {
			return fmt.Errorf("could not create ltv verifier for p7: %w", err)
		}
		if err := l.Verify(); err != nil {
			return fmt.Errorf("verifyLTV information for signature is not valid: %w", err)
		}
		ocspStatus = l.OCSPStatus
	}
	signer := p7.GetOnlySigner()
	signingTime := <-s.signingTime
	if signer.NotBefore.After(signingTime) {
		return fmt.Errorf("certificate was issued at %s, which is after the signing time %s", signer.NotBefore, signingTime)
	}
	signingCertDataResp := signingCertData{
		Signer:      signer.Subject.String(),
		SignerEmail: signer.EmailAddresses[0],
	}
	for _, c := range p7.Certificates {
		cert := CertChain{
			Issuer:    c.Issuer.String(),
			Subject:   c.Subject.String(),
			NotBefore: c.NotBefore,
			NotAfter:  c.NotAfter,
		}
		ocspResponse, ok := ocspStatus[fmt.Sprintf("%x", sha256.Sum256(c.Raw))]
		if ok {
			cert.OCSPStatus = ocspStatusString(ocspResponse.Status)
			cert.OCSPGenerationTime = ocspResponse.ProducedAt
		}

		signingCertDataResp.Certs = append(signingCertDataResp.Certs, cert)
	}
	s.signingCertData <- signingCertDataResp
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

func (s *SignatureContainerVerifier) SendSigningTime(signingTime time.Time) {
	s.signingTime <- signingTime
}

func (s *SignatureContainerVerifier) SigningCertData() signingCertData {
	return <-s.signingCertData
}

func ocspStatusString(status int) string {
	switch status {
	case ocsp.Good:
		return "Good"
	case ocsp.Revoked:
		return "Revoked"
	case ocsp.Unknown:
		return "Unknown"
	default:
		return "ServerFailed"
	}
}