package verifier

import (
	"bytes"
	"crypto"
	"fmt"
)

const (
	ErrHashMismatch = hashError("Hashes didn't match")
)

type hashError string

func (h hashError) Error() string {
	return string(h)
}

func verifyHash(data []byte, hash []byte, algorithm crypto.Hash) error {
	if len(hash) != algorithm.Size() {
		return fmt.Errorf("input hash doesn't size :%d doesn't match expected :%d", len(hash), algorithm.Size())
	}
	hasher := algorithm.New()
	hasher.Write(data)
	if !bytes.Equal(hash, hasher.Sum(nil)) {
		return ErrHashMismatch
	}
	return nil
}
