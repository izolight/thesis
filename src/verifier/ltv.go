package verifier

import (
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"golang.org/x/crypto/ocsp"
)

func VerifyLTV(cert *x509.Certificate, ltvMap map[string]*LTV) error {
	// check if is root CA -> self signed
	if err := cert.CheckSignatureFrom(cert); err == nil {
		return nil
	}
	fingerprint := fmt.Sprintf("%x", sha256.Sum256(cert.Raw))
	ltv, ok := ltvMap[fingerprint]
	if !ok {
		return fmt.Errorf("no ltv information for certificate with fingerprint %s", fingerprint)
	}
	// check first for ocsp and only fallback to crl
	if ltv.Ocsp == nil {
		if ltv.Crl == nil {
			return fmt.Errorf("no ltv information for certificate with fingerprint %s", fingerprint)
		}
	}
	response, err := ocsp.ParseResponseForCert(ltv.Ocsp, cert, nil)
	if err != nil {
		return fmt.Errorf("could not parse ocsp response: %w", err)
	}
	if response.Status != ocsp.Good {
		return fmt.Errorf("certificate with fingerprint %s has ocsp status: %d", fingerprint, response.Status)
	}

	return nil
}