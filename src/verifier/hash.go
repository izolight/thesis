package verifier

import (
	"crypto"
	"errors"
	"fmt"
)

var (
	ErrHashMismatch = errors.New("Hashes didn't match")
)

func verifyHash(data []byte, hash []byte, algorithm crypto.Hash) error {
	if len(hash) != algorithm.Size() {
		return fmt.Errorf("input hash doesn't size :%d doesn't match expected :%d", len(hash), algorithm.Size())
	}
	hasher := algorithm.New()
	hasher.Write(data)
	computedHash := fmt.Sprintf("%x", hasher.Sum(nil))
	inputHash := fmt.Sprintf("%x", hash)
	if computedHash != inputHash {
		return ErrHashMismatch
	}
	return nil
}
