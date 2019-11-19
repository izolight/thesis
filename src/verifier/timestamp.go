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

type TimestampVerifier struct {
	data       chan[]byte
	timestamps []*Timestamped
}

func NewTimestampVerifier(timestamps []*Timestamped) *TimestampVerifier{
	return &TimestampVerifier{
		data:       make(chan []byte, 1),
		timestamps: timestamps,
	}
}

func (t TimestampVerifier) sendData(data []byte) {
	t.data <- data
}

func (t TimestampVerifier) Verify() error {
	if t.timestamps == nil || len(t.timestamps) == 0 {
		return ErrNoTimestamps
	}

	var previousData []byte
	for i, timestamped := range t.timestamps {
		ts, err := pkcs7.ParseTSResponse(timestamped.Rfc3161Timestamp)
		if err != nil {
			return fmt.Errorf("could not parse timestamp response: %w", err)
		}
		l := ltvVerifier{
			certs:  ts.Certificates,
			ltvMap: timestamped.LtvTimestamp,
		}
		err = l.Verify()
		if err != nil {
			return fmt.Errorf("ltv information for timestamp not valid: %w", err)
		}
		hashData := previousData
		// during last signature the data is not in the previous(next) signature,
		// so we need to block until the data arrives
		if i == 0 {
			hashData = <- t.data
		}
		err = verifyHash(hashData, ts.HashedMessage, ts.HashAlgorithm)
		if err != nil {
			return fmt.Errorf("could not verify hash: %w", err)
		}
		previousData, err = proto.Marshal(timestamped)
		if err != nil {
			return fmt.Errorf("could not marshal timestamp: %w", err)
		}
	}
	return nil
}