package verifier

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
)

type signatureDataVerifier struct {
	data         SignatureData
	documentHash []byte
	nonce        chan string
}

func NewSignatureDataVerifier(data SignatureData, documentHash []byte) *signatureDataVerifier {
	return &signatureDataVerifier{
		data:         data,
		documentHash: documentHash,
		nonce:        make(chan string, 1),
	}
}

func (s *signatureDataVerifier) sendNonce(nonce string) {
	s.nonce <- nonce
}

func (s *signatureDataVerifier) Verify() error {
	macAlgo, err := s.data.MacAlgorithm.Algorithm()
	if err != nil {
		return err
	}
	macer := hmac.New(macAlgo.New, s.data.MacKey)
	mac := macer.Sum(s.documentHash)
	foundMAC := false
	for _, m := range s.data.SaltedDocumentHash {
		if bytes.Equal(mac, m) {
			foundMAC = true
			break
		}
	}
	if !foundMAC {
		return errors.New("document hash not found")
	}

	hashAlgo, err := s.data.HashAlgorithm.Algorithm()
	if err != nil {
		return err
	}

	hasher := hashAlgo.New()
	for _, m := range s.data.SaltedDocumentHash {
		hasher.Write(m)
	}
	computedNonce := hasher.Sum(nil)
	nonce, err := hex.DecodeString(<-s.nonce)
	if err != nil {
		return fmt.Errorf("could not decode nonce: %w", err)
	}
	if !bytes.Equal(nonce, computedNonce) {
		return errors.New("computed nonce and id token nonce don't match")
	}

	return nil
}

type macs [][]byte

func (m macs) Len() int {
	return len(m)
}

func (m macs) Less(i, j int) bool {
	return bytes.Compare(m[i], m[j]) <= 0
}

func (m macs) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MACAlgorithm) Algorithm() (crypto.Hash, error) {
	mapping := map[MACAlgorithm]crypto.Hash{
		MACAlgorithm_HMAC_SHA2_256: crypto.SHA256,
		MACAlgorithm_HMAC_SHA2_512: crypto.SHA512,
		MACAlgorithm_HMAC_SHA3_256: crypto.SHA3_256,
		MACAlgorithm_HMAC_SHA3_512: crypto.SHA3_512,
	}
	h, ok := mapping[m]
	if !ok {
		return 0, fmt.Errorf("hash algorithm not implemented :%v", m)
	}
	return h, nil
}

func (h HashAlgorithm) Algorithm() (crypto.Hash, error) {
	return MACAlgorithm(h).Algorithm()
}
