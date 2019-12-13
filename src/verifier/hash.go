package verifier

import (
	"bytes"
	"crypto"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	ErrHashMismatch = hashError("Hashes didn't match")
)

type hashError string

func (h hashError) Error() string {
	return string(h)
}

func verifyHash(data []byte, hash []byte, algorithm crypto.Hash, cfg *Config) error {
	if len(hash) != algorithm.Size() {
		return fmt.Errorf("input hash doesn't size :%d doesn't match expected :%d", len(hash), algorithm.Size())
	}
	hasher := algorithm.New()
	hasher.Write(data)
	calculated_hash := hasher.Sum(nil)
	if !bytes.Equal(hash, calculated_hash) {
		return ErrHashMismatch
	}
	cfg.Logger.WithFields(log.Fields{
		"expected_hash":   fmt.Sprintf("%x", hash),
		"calculated_hash": fmt.Sprintf("%x", calculated_hash),
	}).Info("compared hash")
	return nil
}
