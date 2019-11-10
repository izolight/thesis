package verifier

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type Verifier interface {
	Verify() error
}

func verifySignatureFile(in verifyRequest) error {
	signatureBytes, err := base64.StdEncoding.DecodeString(in.Signature)
	if err != nil {
		return fmt.Errorf("could not decode signature: %w", err)
	}
	signatureFile := &SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		return fmt.Errorf("could not unmarshal signature to protobuf: %w", err)
	}
	data, err := proto.Marshal(signatureFile.SignatureContainer)
	if err != nil {
		return fmt.Errorf("could not marshal signature data: %w", err)
	}
	timestampContainer := timestampContainer{
		data:       data,
		timestamps: signatureFile.GetTimestamps(),
	}
	if err := timestampContainer.Verify(); err != nil {
		return fmt.Errorf("could not verify timestamps: %w", err)
	}
	// TODO: verifySignature -> pkcs#7
	signatureData, err := verifySignature(signatureFile.SignatureContainer)
	if err != nil {
		return fmt.Errorf("could not verify signature: %w", err)
	}
	// TODO: verify id token
	if err := verifyIDToken(signatureData); err != nil {
		return fmt.Errorf("could not verify id token: %w", err)
	}
	// TODO: verify hashes
	if err := verifyHashes(signatureData, in.Hash); err != nil {
		return fmt.Errorf("could not verify hashes: %w", err)
	}

	return nil
}

func verifyHashes(data *SignatureData, hash string) error {
	// TODO: HMAC(salt, hash)
	// TODO: append HMAC to saltedHashes and sort
	// TODO: hash sorted list
	// TODO: compare computed hash with OIDC nonce
	// TODO: compare input hash with signature hash
	return nil
}

func verifyIDToken(data *SignatureData) error {
	// TODO: verify chain and id token
	return nil
}

func verifySignature(container *SignatureContainer) (*SignatureData, error) {
	// TODO: extract data from pkcs7 enveloped data
	var buf []byte
	signatureData := &SignatureData{}
	if err := proto.Unmarshal(buf, signatureData); err != nil {
		return nil, err
	}
	return signatureData, nil
}

