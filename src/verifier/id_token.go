package verifier

type idTokenVerifier struct {
	token []byte
}

func (i idTokenVerifier) Verify() error {
	return nil
}