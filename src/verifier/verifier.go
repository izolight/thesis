package verifier

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
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

func verifySignatureFile(in verifyRequest) error {
	signatureFile, err := decodeSignatureFile(in)
	if err != nil {
		return fmt.Errorf("could not decode signature file: %w", err)
	}

	errors := make(chan error)
	wg := sync.WaitGroup{}

	// TODO add ltvData verifying
	timestampVerifier := NewTimestampVerifier(signatureFile.GetRfc3161InPkcs7(), false, nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := timestampVerifier.Verify(); err != nil {
			errors <- fmt.Errorf("could not verify timestamps: %w", err)
		}
	}()

	timestampVerifier.SendData(signatureFile.SignatureDataInPkcs7)

	signatureContainerVerifier := NewSignatureContainerVerifier(signatureFile.SignatureDataInPkcs7)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := signatureContainerVerifier.Verify(); err != nil {
			errors <- fmt.Errorf("could not verify signatureContainer: %w", err)
		}
	}()

	signatureData := signatureContainerVerifier.getSignatureData()

	signatureDataVerifier := NewSignatureDataVerifier(signatureData, []byte(in.Hash))
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := signatureDataVerifier.Verify(); err != nil {
			errors <- fmt.Errorf("could not verify signatureData: %w", err)
		}
	}()

	idTokenVerifier, err := NewIDTokenVerifier(&signatureData, &cfg, timestampVerifier.getNotAfter(), false)
	if err != nil {
		return fmt.Errorf("could not create id token verifier: %w", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := idTokenVerifier.Verify(); err != nil {
			errors <- fmt.Errorf("could not verify id token: %w", err)
		}
	}()
	signatureDataVerifier.sendNonce(idTokenVerifier.getNonce())

	wg.Wait()

	err = <-errors
	if err != nil {
		return fmt.Errorf("error during verification: %w", err)
	}

	return nil
}

func decodeSignatureFile(in verifyRequest) (*SignatureFile, error) {
	signatureBytes, err := base64.StdEncoding.DecodeString(in.Signature)
	if err != nil {
		return nil, fmt.Errorf("could not decode signature: %w", err)
	}
	signatureFile := &SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		return nil, fmt.Errorf("could not unmarshal signature to protobuf: %w", err)
	}
	return signatureFile, nil
}
