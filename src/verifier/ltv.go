package verifier

import (
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"golang.org/x/crypto/ocsp"
)

type ltvVerifier struct {
	certs []*x509.Certificate
	ltvMap map[string]*LTV
}

func (l ltvVerifier) Verify() error {
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
		ltv, ok := l.ltvMap[fingerprint]
		if !ok {
			return fmt.Errorf("no ltv information for certificate with fingerprint %s", fingerprint)
		}
		// check first for ocsp and only fallback to crl
		if ltv.Ocsp == nil {
			if ltv.Crl == nil {
				return fmt.Errorf("no ltv information for certificate with fingerprint %s", fingerprint)
			}
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