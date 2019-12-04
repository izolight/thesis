package verifier

import (
	"fmt"
	"sync"
)

type Verifier interface {
	Verify() error
}

type Config struct {
	Issuer   string
	ClientId string
}

var (
	cfg = Config{
		Issuer:   "",
		ClientId: "",
	}
)

func VerifySignatureFile(file *SignatureFile, hash string) error {
	errors := make(chan error)
	wg := sync.WaitGroup{}

	// TODO add ltvData verifying
	timestampVerifier := NewTimestampVerifier(file.GetRfc3161InPkcs7(), false, nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := timestampVerifier.Verify(); err != nil {
			errors <- fmt.Errorf("could not verify timestamps: %w", err)
		}
	}()

	timestampVerifier.SendData(file.SignatureDataInPkcs7)

	signatureContainerVerifier := NewSignatureContainerVerifier(file.SignatureDataInPkcs7)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := signatureContainerVerifier.Verify(); err != nil {
			errors <- fmt.Errorf("could not verify signatureContainer: %w", err)
		}
	}()

	signatureData := signatureContainerVerifier.getSignatureData()

	signatureDataVerifier, err := NewSignatureDataVerifier(&signatureData, hash)
	if err != nil {
		return fmt.Errorf("could not create signature data verifier: %w", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := signatureDataVerifier.Verify(false); err != nil {
			errors <- fmt.Errorf("could not verify signatureData: %w", err)
		}
	}()

	idTokenVerifier, err := NewIDTokenVerifier(&signatureData, &cfg, timestampVerifier.getNotAfter())
	if err != nil {
		return fmt.Errorf("could not create id token verifier: %w", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := idTokenVerifier.Verify(false); err != nil {
			errors <- fmt.Errorf("could not verify id token: %w", err)
		}
	}()
	signatureDataVerifier.SendNonce(idTokenVerifier.getNonce())

	wg.Wait()

	err = <-errors
	if err != nil {
		return fmt.Errorf("error during verification: %w", err)
	}

	return nil
}
