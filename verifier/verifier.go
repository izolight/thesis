package verifier

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

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
		signatureDataVerifier := NewSignatureDataVerifier(&signatureData, hash, s.cfg)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := signatureDataVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify signatureData: %w", err)
			}
		}()

		timestampData := timestampVerifier.TimestampData()
		signatureContainerVerifier.SendSigningTime(timestampData.SigningTime)
		s.cfg.Logger.WithFields(log.Fields{
			"signing_time": timestampData.SigningTime,
		}).Info("decoded signing time")
		idTokenVerifier, err := NewIDTokenVerifier(&signatureData, timestampData.SigningTime, s.cfg)
		if err != nil {
			errors <- fmt.Errorf("could not create id token verifier: %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := idTokenVerifier.Verify(false); err != nil {
				errors <- fmt.Errorf("could not verify id token: %w", err)
			}
		}()

		idToken := idTokenVerifier.IDToken()
		signatureDataVerifier.SendNonce(idToken.Nonce)
		signingCertData := signatureContainerVerifier.SigningCertData()
		idTokenVerifier.SendEmail(signingCertData.SignerEmail)

		wg.Wait()
		resp := VerifyResponse{
			Valid:       true,
			Error:       "",
			Signature:   signatureDataVerifier.SignatureData(),
			SigningCert: signingCertData,
			Timestamp:   timestampData,
			IDToken:     idToken,
		}

		responses <- resp
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

type VerifyResponse struct {
	Valid       bool              `json:"valid"`
	Error       string            `json:"error,omitempty"`
	IDToken     idToken           `json:"id_token"`
	Signature   signatureData     `json:"signature"`
	SigningCert signingCertData   `json:"signing_cert"`
	Timestamp   timestampDataResp `json:"timestamp"`
}

type CertChain struct {
	Issuer    string    `json:"issuer"`
	Subject   string    `json:"subject"`
	NotBefore time.Time `json:"not_before"`
	NotAfter  time.Time `json:"not_after"`
}