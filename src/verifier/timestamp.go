package verifier

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mozilla.org/pkcs7"
	"time"
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
	notAfter chan time.Time
	timestamps []*Timestamped
}

func NewTimestampVerifier(timestamps []*Timestamped) *TimestampVerifier{
	return &TimestampVerifier{
		data:       make(chan []byte, 1),
		notAfter: make(chan time.Time, 1),
		timestamps: timestamps,
	}
}

func (t *TimestampVerifier) sendData(data []byte) {
	t.data <- data
}

func (t *TimestampVerifier) getNotAfter() time.Time {
	return <- t.notAfter
}

func verifyTimestamp(t *Timestamped, data []byte) (*time.Time, error) {
	ts, err := pkcs7.ParseTSResponse(t.Rfc3161Timestamp)
	if err != nil {
		return nil, fmt.Errorf("could not parse timestamp response: %w", err)
	}
	l := ltvVerifier{
		certs:  ts.Certificates,
		ltvMap: t.LtvTimestamp,
	}

	if err = l.Verify(); err != nil {
		return nil, fmt.Errorf("ltv information for timestamp not valid: %w", err)
	}
	err = verifyHash(data, ts.HashedMessage, ts.HashAlgorithm)
	if err != nil {
		return nil, fmt.Errorf("could not verify hash: %w", err)
	}
	return &ts.Time, nil
}

func (t *TimestampVerifier) Verify() error {
	if t.timestamps == nil || len(t.timestamps) == 0 {
		return ErrNoTimestamps
	}

	for i := len(t.timestamps)-1; i >= 0; i-- {
		hashData := []byte{}
		// during last signature the data is not in the previous(next) signature,
		// so we need to block until the data arrives
		if i == 0 {
			hashData = <- t.data
		} else {
			var err error
			hashData, err = proto.Marshal(t.timestamps[i-1])
			if err != nil {
				return fmt.Errorf("could not marshal timestamp: %w", err)
			}
		}
		notAfter, err := verifyTimestamp(t.timestamps[i], hashData)
		if err != nil {
			return fmt.Errorf("could not verify timestamp: %w", err)
		}
		if i == 0 {
			t.notAfter <- *notAfter
		}
	}
	return nil
}