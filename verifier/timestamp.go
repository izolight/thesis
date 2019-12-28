package verifier

import (
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
	data          []byte
	timestamps    [][]byte
	verifyLTV     bool
	timestampData chan timestampDataResp
	ltvData       map[string]*LTV
	cfg           *Config
}

type timestampDataResp struct {
	SigningTime time.Time
	Certs       []CertChain `json:"cert_chain"`
}

func NewTimestampVerifier(timestamps [][]byte, data []byte, verifyLTV bool, ltvData map[string]*LTV, cfg Config) *TimestampVerifier {
	cfg.Logger = cfg.Logger.WithField("verifier", "timestamp")
	return &TimestampVerifier{
		data:          data,
		timestamps:    timestamps,
		verifyLTV:     verifyLTV,
		ltvData:       ltvData,
		cfg:           &cfg,
		timestampData: make(chan timestampDataResp, 1),
	}
}

func (t *TimestampVerifier) verifyTimestamp(timestamp []byte, data []byte, index int) error {
	ts, err := pkcs7.ParseTSResponse(timestamp)
	if err != nil {
		return err
	}
	if t.verifyLTV {
		l := LTVVerifier{
			certs: ts.Certificates,
			//LTVData: t.ltvData,
		}
		if err := l.Verify(); err != nil {
			return err
		}
	}

	if err = verifyHash(data, ts.HashedMessage, ts.HashAlgorithm, *t.cfg); err != nil {
		return err
	}
	if index == 0 {
		timestampData := timestampDataResp{
			SigningTime: ts.Time,
		}
		for _, c := range ts.Certificates {
			timestampData.Certs = append(timestampData.Certs, CertChain{
				Issuer:    c.Issuer.String(),
				Subject:   c.Subject.String(),
				NotBefore: c.NotBefore,
				NotAfter:  c.NotAfter,
			})
		}
		t.timestampData <- timestampData
	}
	t.cfg.Logger.WithFields(log.Fields{
		"timestamp":        ts.Time,
		"timestamped_hash": fmt.Sprintf("%x", ts.HashedMessage),
	}).Info("verified timestamp")

	return nil
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
		if err := t.verifyTimestamp(t.timestamps[i], hashData, i); err != nil {
			return fmt.Errorf("could not verify timestamp: %w", err)
		}
	}
	t.cfg.Logger.Info("finished verifying")
	return nil
}

func (t *TimestampVerifier) TimestampData() timestampDataResp {
	return <-t.timestampData
}
