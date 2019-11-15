package verifier

type signatureDataVerifier struct {
	data SignatureData
}

func (s signatureDataVerifier) Verify() error {
	return nil
}