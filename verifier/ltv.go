package verifier

import (
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"golang.org/x/crypto/ocsp"
)

type LTVVerifier struct {
	certs         []*x509.Certificate
	crls          []pkix.CertificateList
	ocspResponses [][]byte

	OCSPStatus map[string]*ocsp.Response
	CRLStatus map[string]*pkix.CertificateList
}

func NewLTVVerifier(certs []*x509.Certificate, crls []pkix.CertificateList, ocsps [][]byte) *LTVVerifier {
	return &LTVVerifier{
		certs:         certs,
		crls:          crls,
		ocspResponses: ocsps,
		OCSPStatus:    make(map[string]*ocsp.Response),
		CRLStatus:     make(map[string]*pkix.CertificateList),
	}
}

func (l LTVVerifier) Verify() error {
	for _, cert := range l.certs {
		// check if is root CA -> no ocsp/crl possible
		if err := cert.CheckSignatureFrom(cert); err == nil {
			continue
		}
		var issuingCA *x509.Certificate
		for _, issuing := range l.certs {
			if err := cert.CheckSignatureFrom(issuing); err == nil {
				issuingCA = issuing
				break
			}
		}

		fingerprint := fmt.Sprintf("%x", sha256.Sum256(cert.Raw))

		// check first for ocsp and only fallback to crl
		if len(l.ocspResponses) == 0 {
			if len(l.crls) == 0 {
				return errors.New("no crls found")
			}
			return errors.New("crl not supported yet")
		}


		for _, ocspResponse := range l.ocspResponses {
			response, err := ocsp.ParseResponseForCert(ocspResponse, cert, issuingCA)
			if err != nil {
				return fmt.Errorf("couldn't verify ocsp response for %s: %w", cert.Subject.String(), err)
			}
			if response.Status != ocsp.Good {
				return fmt.Errorf("certificate %s has ocsp status: %d", cert.Subject.String(), response.Status)
			}
			l.OCSPStatus[fingerprint] = response
			break
		}
	}
	return nil
}
