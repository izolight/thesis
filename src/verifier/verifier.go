package verifier

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mozilla.org/pkcs7"
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
	var data []byte
	_, err = signatureFile.SignatureData.XXX_Marshal(data, true)
	if err != nil {
		return false, fmt.Errorf("could not marshal signature data: %w", err)
	}
	if err := verifyTimestamps(data, signatureFile.GetTimestamps()); err != nil {
		return false, fmt.Errorf("could not verify timestamps: %w", err)
	}
	return true, nil
}

func verifyTimestamps(data []byte, timestamps []*Timestamped) error {
	if timestamps == nil {
		return errors.New("No timestamps included")
	}

	var previousBytes []byte
	for i, timestamped := range timestamps {
		ts, err := pkcs7.ParseTSResponse(timestamped.Rfc3161Timestamp)
		if err != nil {
			return fmt.Errorf("could not parse timestamp response: %w", err)
		}
		hashData := previousBytes
		if i == 0 {
			hashData = data
		}
		hasher := ts.HashAlgorithm.New()
		hasher.Write(hashData)
		hash := fmt.Sprintf("%x", hasher.Sum(nil))
		tsHash := fmt.Sprintf("%x", ts.HashedMessage)
		if hash != tsHash {
			return fmt.Errorf("timestamped hashes didn't match: %s != %s", hash, tsHash)
		}
		previousBytes = timestamped.Rfc3161Timestamp
	}

	return nil
}