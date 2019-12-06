package verifier

import (
	"fmt"
	"sync"
)

type Config struct {
	Issuer   string
	ClientId string
}

type SignatureVerifier struct {
	cfg Config
}

func NewSignatureVerifier(cfg Config) *SignatureVerifier {
	return &SignatureVerifier{cfg:cfg}
}

func (s SignatureVerifier) VerifySignatureFile(file *SignatureFile, hash string) (verifyResponse, error) {
	errors := make(chan error, 1)
	responses := make(chan verifyResponse)
	wg := sync.WaitGroup{}

	// TODO add ltvData verifying
	go func() {
		timestampVerifier := NewTimestampVerifier(file.GetRfc3161InPkcs7(), file.GetSignatureDataInPkcs7(), false, nil)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := timestampVerifier.Verify(); err != nil {
				errors <- fmt.Errorf("could not verify timestamps: %w", err)
			}
		}()

		signatureContainerVerifier := NewSignatureContainerVerifier(file.SignatureDataInPkcs7)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := signatureContainerVerifier.Verify(); err != nil {
				errors <- fmt.Errorf("could not verify signatureContainer: %w", err)
			}

		}()

		signatureData := signatureContainerVerifier.SignatureData()
		signatureDataVerifier, err := NewSignatureDataVerifier(&signatureData, hash)
		if err != nil {
			errors <- fmt.Errorf("could not create signature data verifier: %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := signatureDataVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify signatureData: %w", err)
			}
		}()

		signingTime := timestampVerifier.SigningTime()
		idTokenVerifier, err := NewIDTokenVerifier(&signatureData, &s.cfg, signingTime)
		if err != nil {
			errors <- fmt.Errorf("could not create id token verifier: %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := idTokenVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify id token: %w", err)
			}
			signatureDataVerifier.SendNonce(idTokenVerifier.getNonce())
		}()


		wg.Wait()
		responses <- verifyResponse{
			Valid:          true,
			Error:          "",
			SignerEmail:    signatureContainerVerifier.SignerEmail(),
			SignatureLevel: signatureData.SignatureLevel,
			SignatureTime:  signingTime,
		}
	}()

	var err error
	var resp verifyResponse
	select {
	case err = <-errors:
		break
	case resp = <-responses:
		break
	}
	return resp, err
}
