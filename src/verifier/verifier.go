package verifier

import (
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	Issuer          string
	ClientId        string
	AdditionalCerts []*x509.Certificate
	Logger          *log.Entry
}

type SignatureVerifier struct {
	cfg Config
}

func NewSignatureVerifier(cfg Config) *SignatureVerifier {
	return &SignatureVerifier{cfg: cfg}
}

func (s *SignatureVerifier) VerifySignatureFile(file *SignatureFile, hash string) (VerifyResponse, error) {
	errors := make(chan error, 1)
	responses := make(chan VerifyResponse)
	wg := sync.WaitGroup{}

	// TODO add ltvData verifying
	go func() {
		timestampVerifier := NewTimestampVerifier(file.GetRfc3161InPkcs7(), file.GetSignatureDataInPkcs7(), false, nil, s.cfg)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := timestampVerifier.Verify(); err != nil {
				errors <- fmt.Errorf("could not verify timestamps: %w", err)
			}
		}()

		signatureContainerVerifier := NewSignatureContainerVerifier(file.SignatureDataInPkcs7, s.cfg.AdditionalCerts, s.cfg)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := signatureContainerVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify signatureContainer: %w", err)
			}
		}()

		signatureData := signatureContainerVerifier.SignatureData()
		signatureDataVerifier, err := NewSignatureDataVerifier(&signatureData, hash, s.cfg)
		if err != nil {
			errors <- fmt.Errorf("could not create signature data verifier: %w", err)
		}

		wg.Add(1)
		go func(logger *log.Entry) {
			defer wg.Done()
			if err := signatureDataVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify signatureData: %w", err)
			}
		}(s.cfg.Logger)

		signingTime := timestampVerifier.SigningTime()
		signatureContainerVerifier.SendSigningTime(signingTime)
		s.cfg.Logger.WithFields(log.Fields{
			"signing_time": signingTime,
		}).Info("decoded signing time")
		idTokenVerifier, err := NewIDTokenVerifier(&signatureData, signingTime, s.cfg)
		if err != nil {
			errors <- fmt.Errorf("could not create id token verifier: %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := idTokenVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify id token: %w", err)
			}
			signatureDataVerifier.SendNonce(idTokenVerifier.Nonce())
		}()
		signer := signatureContainerVerifier.Signer()
		idTokenVerifier.SendSigner(signer)

		wg.Wait()
		responses <- VerifyResponse{
			Valid:          true,
			Error:          "",
			SignerEmail:    signer.EmailAddresses[0],
			SignatureLevel: signatureData.SignatureLevel,
			SignatureTime:  signingTime,
		}
	}()

	var err error
	var resp VerifyResponse
	select {
	case err = <-errors:
		break
	case resp = <-responses:
		break
	}
	return resp, err
}
