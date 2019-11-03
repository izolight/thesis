package verifier

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func verifySignature(in verifyRequest) (bool, error) {
	signatureBytes, err := base64.StdEncoding.DecodeString(in.Signature)
	if err != nil {
		return false, fmt.Errorf("could not decode signature: %w", err)
	}
	signatureFile := &SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		return false, fmt.Errorf("could not unmarshal signature to protobuf: %w", err)
	}
	if err := verifyTimestamps(signatureFile.GetTimestamps()); err != nil {
		return false, fmt.Errorf("could not verify timestamps: %w", err)
	}
	return true, nil
}

func verifyTimestamps(timestamps []*Timestamped) error {
	if timestamps == nil {
		return errors.New("No timestamps included")
	}

	return nil
}