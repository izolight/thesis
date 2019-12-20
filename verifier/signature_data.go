package verifier

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type signatureDataVerifier struct {
	data         *SignatureData
	documentHash string
	nonce        chan string
	signatureData chan signatureDataResp
	cfg          *Config
}

type signatureDataResp struct {
	SaltedHashes []string `json:"salted_hashes"`
	HashAlgorithm string `json:"hash_algorithm"`
	MacKey string `json:"mac_key"`
	MACAlgorithm string `json:"mac_algorithm"`
	SignatureLevel string `json:"signature_level"`
}

func NewSignatureDataVerifier(data *SignatureData, documentHash string, cfg Config) *signatureDataVerifier {
	cfg.Logger = cfg.Logger.WithField("verifier", "signature data")
	v := &signatureDataVerifier{
		data:  data,
		nonce: make(chan string, 1),
		signatureData: make(chan signatureDataResp, 1),
		cfg:   &cfg,
		documentHash: documentHash,
	}
	return v
}

func (s *signatureDataVerifier) SendNonce(nonce string) {
	s.nonce <- nonce
}

func (s *signatureDataVerifier) SignatureData() signatureDataResp {
	return <- s.signatureData
}

func (s *signatureDataVerifier) Verify(verifyLTV bool) error {
	s.cfg.Logger.Info("started verifying")
	macAlgo, err := s.data.MacAlgorithm.Algorithm()
	if err != nil {
		return err
	}
	macer := hmac.New(macAlgo.New, s.data.MacKey)
	hashBytes, err := hex.DecodeString(s.documentHash)
	if err != nil {
		return fmt.Errorf("could not decode document hash: %w", err)
	}
	macer.Write(hashBytes)
	mac := macer.Sum(nil)
	s.cfg.Logger.WithFields(log.Fields{
		"mac":           fmt.Sprintf("%x", mac),
		"mac_key":       fmt.Sprintf("%x", s.data.MacKey),
		"document_hash": s.documentHash,
		"mac_algorithm": s.data.MacAlgorithm,
	}).Info("calculated mac")
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
	var saltedHashes []string
	for i := range s.data.SaltedDocumentHash {
		saltedHashes = append(saltedHashes, fmt.Sprintf("%x", s.data.SaltedDocumentHash[i]))
	}
	s.signatureData <- signatureDataResp{
		SaltedHashes:   saltedHashes,
		HashAlgorithm:  s.data.HashAlgorithm.String(),
		MacKey:         fmt.Sprintf("%x", s.data.MacKey),
		MACAlgorithm:   s.data.MacAlgorithm.String(),
		SignatureLevel: s.data.SignatureLevel.String(),
	}

	s.cfg.Logger.WithFields(log.Fields{
		"hash_algorithm":         s.data.HashAlgorithm,
		"salted_document_hashes": fmt.Sprintf("%x", s.data.SaltedDocumentHash),
		"computed_nonce":         fmt.Sprintf("%x", computedNonce),
		"id_token_nonce":         fmt.Sprintf("%x", nonce),
	}).Info("calculated nonce")

	s.cfg.Logger.Info("finished verifying")
	return nil
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
