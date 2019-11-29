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

type config struct {
	issuer string
	clientId string
}

var (
	cfg = config{
		issuer:   "",
		clientId: "",
	}
)

func verifySignatureFile(in verifyRequest) error {
	signatureFile, err := decodeSignatureFile(in)
	if err != nil {
		return fmt.Errorf("could not decode signature file: %w", err)
	}

	errors := make(chan error)
	wg := sync.WaitGroup{}

	timestampVerifier := NewTimestampVerifier(signatureFile.GetTimestamps())
	go func() {
		wg.Add(1)
		err := timestampVerifier.Verify()
		if err != nil {
			errors <- fmt.Errorf("could not verify timestamps: %w", err)
		}
		wg.Done()
	}()

	data, err := proto.Marshal(signatureFile.SignatureContainer)
	if err != nil {
		return fmt.Errorf("could not marshal signature Data: %w", err)
	}
	timestampVerifier.sendData(data)

	signatureContainerVerifier := NewSignatureContainerVerifier(signatureFile.SignatureContainer)
	go func() {
		wg.Add(1)
		err := signatureContainerVerifier.Verify()
		if err != nil {
			errors <- fmt.Errorf("could not verify signatureContainer: %w", err)
		}
		wg.Done()
	}()

	signatureData := signatureContainerVerifier.getSignatureData()

	signatureDataVerifier := NewSignatureDataVerifier(signatureData, []byte(in.Hash))
	go func() {
		wg.Add(1)
		err := signatureDataVerifier.Verify()
		if err != nil {
			errors <- fmt.Errorf("could not verify signatureData: %w", err)
		}
		wg.Done()
	}()

	idTokenVerifier, err := NewIDTokenVerifier(&signatureData, &cfg, timestampVerifier.getNotAfter())
	if err != nil {
		return fmt.Errorf("could not create id token verifier: %w", err)
	}
	go func() {
		wg.Add(1)
		err := idTokenVerifier.Verify()
		if err != nil {
			errors <- fmt.Errorf("could not verify id token: %w", err)
		}
		wg.Done()
	}()
	signatureDataVerifier.sendNonce(idTokenVerifier.getNonce())

	wg.Wait()

	err = <- errors
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