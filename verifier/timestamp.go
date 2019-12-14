package verifier

import (
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	timestamps  [][]byte
	verifyLTV   bool
	signingTime chan time.Time
	certs chan []*x509.Certificate
	ltvData     map[string]*LTV
	cfg         *Config
}

func NewTimestampVerifier(timestamps [][]byte, data []byte, verifyLTV bool, ltvData map[string]*LTV, cfg Config) *TimestampVerifier {
	cfg.Logger = cfg.Logger.WithField("verifier", "timestamp")
	return &TimestampVerifier{
		data:        data,
		timestamps:  timestamps,
		verifyLTV:   verifyLTV,
		ltvData:     ltvData,
		cfg:         &cfg,
		signingTime: make(chan time.Time, 1),
		certs: make(chan []*x509.Certificate, 1),
	}
}

func (t *TimestampVerifier) verifyTimestamp(timestamp []byte, data []byte, index int) (*time.Time, error) {
	ts, err := pkcs7.ParseTSResponse(timestamp)
	if err != nil {
		return nil, err
	}
	if t.verifyLTV {
		l := LTVVerifier{
			Certs:   ts.Certificates,
			LTVData: t.ltvData,
		}
		if err := l.Verify(); err != nil {
			return nil, err
		}
	}

	if err = verifyHash(data, ts.HashedMessage, ts.HashAlgorithm, *t.cfg); err != nil {
		return nil, err
	}
	if index == 0 {
		t.certs <- ts.Certificates
	}
	t.cfg.Logger.WithFields(log.Fields{
		"timestamp":        ts.Time,
		"timestamped_hash": fmt.Sprintf("%x", ts.HashedMessage),
	}).Info("verified timestamp")

	return &ts.Time, nil
}

func (t *TimestampVerifier) Verify() error {
	t.cfg.Logger.Info("started verifying")
	if t.timestamps == nil || len(t.timestamps) == 0 {
		return ErrNoTimestamps
	}

	t.cfg.Logger.Infof("got %d timestamps", len(t.timestamps))
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
		notAfter, err := t.verifyTimestamp(t.timestamps[i], hashData, i)
		if err != nil {
			return fmt.Errorf("could not verify timestamp: %w", err)
		}
		if i == 0 {
			t.signingTime <- *notAfter
		}
	}
	t.cfg.Logger.Info("finished verifying")
	return nil
}

func (t *TimestampVerifier) SigningTime() time.Time {
	return <-t.signingTime
}

func (t *TimestampVerifier) Certs() []*x509.Certificate {
	return <-t.certs
}