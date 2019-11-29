package verifier

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"errors"
	"fmt"
	"sort"
)

type signatureDataVerifier struct {
	data SignatureData
	documentHash []byte
	nonce chan string
}

func NewSignatureDataVerifier(data SignatureData, documentHash []byte) *signatureDataVerifier{
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
	if !bytes.Equal(s.documentHash, s.data.DocumentHash) {
		return fmt.Errorf("document hash and signature didn't match")
	}
	macAlgo,err := s.data.MacAlgorithm.Algorithm()
	if err != nil {
		return err
	}
	macer := hmac.New(macAlgo.New, s.data.MacKey)
	mac := macer.Sum(s.documentHash)
	allMacs := append(s.data.OtherMacs, mac)
	sort.Sort(macs(allMacs))

	hashAlgo,err := s.data.HashAlgorithm.Algorithm()
	if err != nil {
		return err
	}

	hasher := hashAlgo.New()
	for i := range allMacs {
		hasher.Write(allMacs[i])
	}
	computedNonce := hasher.Sum(nil)
	nonce := []byte(<-s.nonce)
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