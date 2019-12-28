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
	Certs   []*x509.Certificate
	LTVData map[string]*LTV
}

func NewLTVVerifier(certs []*x509.Certificate, crls []pkix.CertificateList, ocsps []ocsp.Response) (*LTVVerifier, error) {
	l := &LTVVerifier{
		Certs: certs,
		LTVData: make(map[string]*LTV),
	}
	for _, o := range ocsps {
		l.LTVData[fmt.Sprintf("%x", o.Certificate.Raw)] = &LTV{Ocsp:o.TBSResponseData}
	}

	return l, nil
}

func (l LTVVerifier) Verify() error {
	for _, cert := range l.Certs {
		// check if is root CA -> no ocsp/crl possible
		if err := cert.CheckSignatureFrom(cert); err == nil {
			continue
		}
		var issuingCA *x509.Certificate
		for _, issuing := range l.Certs {
			if err := cert.CheckSignatureFrom(issuing); err == nil {
				issuingCA = issuing
				break
			}
		}

		fingerprint := fmt.Sprintf("%x", sha256.Sum256(cert.Raw))
		ltv, ok := l.LTVData[fingerprint]
		if !ok || ltv == nil {
			return fmt.Errorf("no verifyLTV information for certificate with fingerprint %s", fingerprint)
		}
		// check first for ocsp and only fallback to crl
		if ltv.Ocsp == nil {
			if ltv.Crl == nil {
				return fmt.Errorf("no verifyLTV information for certificate with fingerprint %s", fingerprint)
			}
			return errors.New("crl not supported yet")
			// TODO: verify crl
		}
		response, err := ocsp.ParseResponseForCert(ltv.Ocsp, cert, issuingCA)
		if err != nil {
			return fmt.Errorf("could not parse ocsp response: %w", err)
		}
		if response.Status != ocsp.Good {
			return fmt.Errorf("certificate with fingerprint %s has ocsp status: %d", fingerprint, response.Status)
		}
	}
	return nil
}
