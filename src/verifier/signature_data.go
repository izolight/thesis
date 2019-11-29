package verifier

import "errors"

type signatureDataVerifier struct {
	data SignatureData
	hmac []byte
}

func (s signatureDataVerifier) Verify() error {

	return errors.New("not implemented")
}