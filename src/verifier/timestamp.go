package verifier

import (
	"fmt"
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
	data        []byte
	signingTime chan (time.Time)
	timestamps  [][]byte
	verifyLTV   bool
	ltvData     map[string]*LTV
}

func NewTimestampVerifier(timestamps [][]byte, data []byte, verifyLTV bool, ltvData map[string]*LTV) *TimestampVerifier {
	return &TimestampVerifier{
		data:        data,
		signingTime: make(chan time.Time, 1),
		timestamps:  timestamps,
		verifyLTV:   verifyLTV,
		ltvData:     ltvData,
	}
}

func (t *TimestampVerifier) SigningTime() time.Time {
	return <-t.signingTime
}

func (t *TimestampVerifier) verifyTimestamp(timestamp []byte, data []byte) (*time.Time, error) {
	ts, err := pkcs7.ParseTSResponse(timestamp)
	if err != nil {
		return nil, fmt.Errorf("could not parse timestamp response: %w", err)
	}
	if t.verifyLTV {
		l := LTVVerifier{
			Certs:   ts.Certificates,
			LTVData: t.ltvData,
		}
		if err = l.Verify(); err != nil {
			return nil, fmt.Errorf("verifyLTV information for timestamp not valid: %w", err)
		}
	}

	if err = verifyHash(data, ts.HashedMessage, ts.HashAlgorithm); err != nil {
		return nil, fmt.Errorf("could not verify hash: %w", err)
	}
	return &ts.Time, nil
}

func (t *TimestampVerifier) Verify() error {
	if t.timestamps == nil || len(t.timestamps) == 0 {
		return ErrNoTimestamps
	}

	for i := len(t.timestamps) - 1; i >= 0; i-- {
		var hashData []byte
		// during last signature the data is not in the previous(next) signature,
		// so we need to block until the data arrives
		if i == 0 {
			hashData = t.data
		} else {
			hashData = t.timestamps[i-1]
			//var err error
			//hashData, err = proto.Marshal(t.timestamps[i-1])
			//if err != nil {
			//	return fmt.Errorf("could not marshal timestamp: %w", err)
			//}
		}
		notAfter, err := t.verifyTimestamp(t.timestamps[i], hashData)
		if err != nil {
			return fmt.Errorf("could not verify timestamp: %w", err)
		}
		if i == 0 {
			t.signingTime <- *notAfter
		}
	}
	return nil
}
