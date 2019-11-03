package verifier

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mozilla.org/pkcs7"
)

func verifySignatureFile(in verifyRequest) error {
	signatureBytes, err := base64.StdEncoding.DecodeString(in.Signature)
	if err != nil {
		return fmt.Errorf("could not decode signature: %w", err)
	}
	signatureFile := &SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		return fmt.Errorf("could not unmarshal signature to protobuf: %w", err)
	}
	var data []byte
	_, err = signatureFile.SignatureContainer.XXX_Marshal(data, true)
	if err != nil {
		return fmt.Errorf("could not marshal signature data: %w", err)
	}
	if err := verifyTimestamps(data, signatureFile.GetTimestamps()); err != nil {
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
		// TODO: verify ocsp and crl for each timestamp
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