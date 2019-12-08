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
		timestampVerifier := NewTimestampVerifier(file.GetRfc3161InPkcs7(), file.GetSignatureDataInPkcs7(), false, nil)
		wg.Add(1)
		go func(logger *log.Entry) {
			defer wg.Done()
			logger.Info("start timestamp verifying")
			if err := timestampVerifier.Verify(); err != nil {
				errors <- fmt.Errorf("could not verify timestamps: %w", err)
			}
			logger.Info("finished timestamp verifying")
		}(s.cfg.Logger)

		signatureContainerVerifier := NewSignatureContainerVerifier(file.SignatureDataInPkcs7, s.cfg.AdditionalCerts)
		wg.Add(1)
		go func(logger *log.Entry) {
			defer wg.Done()
			logger.Info("start signature container verifying")
			if err := signatureContainerVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify signatureContainer: %w", err)
			}
			logger.Info("finished signature container verifying")
		}(s.cfg.Logger)

		signatureData := signatureContainerVerifier.SignatureData()
		s.cfg.Logger.WithFields(log.Fields{
			"signatureData": signatureData.String(),
		}).Info()
		signatureDataVerifier, err := NewSignatureDataVerifier(&signatureData, hash)
		if err != nil {
			errors <- fmt.Errorf("could not create signature data verifier: %w", err)
		}

		wg.Add(1)
		go func(logger *log.Entry) {
			defer wg.Done()
			logger.Info("start signature data verifying")
			if err := signatureDataVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify signatureData: %w", err)
			}
			logger.Info("finished signature data verifying")
		}(s.cfg.Logger)

		signingTime := timestampVerifier.SigningTime()
		idTokenVerifier, err := NewIDTokenVerifier(&signatureData, &s.cfg, signingTime)
		if err != nil {
			errors <- fmt.Errorf("could not create id token verifier: %w", err)
		}

		wg.Add(1)
		go func(logger *log.Entry) {
			defer wg.Done()
			logger.Info("start id token verifying")
			if err := idTokenVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify id token: %w", err)
			}
			logger.Info("finished id token verifying")
			signatureDataVerifier.SendNonce(idTokenVerifier.getNonce())
		}(s.cfg.Logger)

		wg.Wait()
		responses <- VerifyResponse{
			Valid:          true,
			Error:          "",
			SignerEmail:    signatureContainerVerifier.SignerEmail(),
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
