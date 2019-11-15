package verifier

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mozilla.org/pkcs7"
)

const (
	ErrNoTimestamps = timestampError("No timestamps included")
)

type timestampError string

func (t timestampError) Error() string {
	return string(t)
}

type timestampVerifier struct {
	data []byte
	timestamps []*Timestamped
}

func (t timestampVerifier) Verify() error {
	if t.timestamps == nil || len(t.timestamps) == 0 {
		return ErrNoTimestamps
	}

	var previousBytes []byte
	for i, timestamped := range t.timestamps {
		ts, err := pkcs7.ParseTSResponse(timestamped.Rfc3161Timestamp)
		if err != nil {
			return fmt.Errorf("could not parse timestamp response: %w", err)
		}
		// TODO: verify ocsp and crl for each timestamp
		hashData := previousBytes
		if i == 0 {
			hashData = t.data
		}
		hasher := ts.HashAlgorithm.New()
		hasher.Write(hashData)
		hash := fmt.Sprintf("%x", hasher.Sum(nil))
		tsHash := fmt.Sprintf("%x", ts.HashedMessage)
		if hash != tsHash {
			return fmt.Errorf("timestamped hashes didn't match: %s != %s", hash, tsHash)
		}
		previousBytes, err = proto.Marshal(timestamped)
		if err != nil {
			return fmt.Errorf("could not marshal timestamp: %w", err)
		}
	}
	return nil
}