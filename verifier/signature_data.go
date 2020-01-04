package verifier

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/pb"
)

type signatureDataVerifier struct {
	data          *pb.SignatureData
	documentHash  string
	nonce         chan string
	signatureData chan signatureData
	cfg           *Config
}

type signatureData struct {
	SaltedHashes   []string `json:"salted_hashes"`
	HashAlgorithm  string   `json:"hash_algorithm"`
	MacKey         string   `json:"mac_key"`
	MACAlgorithm   string   `json:"mac_algorithm"`
	SignatureLevel string   `json:"signature_level"`
}

func NewSignatureDataVerifier(data *pb.SignatureData, documentHash string, cfg Config) *signatureDataVerifier {
	cfg.Logger = cfg.Logger.WithField("verifier", "signature data")
	v := &signatureDataVerifier{
		data:          data,
		nonce:         make(chan string, 1),
		signatureData: make(chan signatureData, 1),
		cfg:           &cfg,
		documentHash:  documentHash,
	}
	return v
}

func (s *signatureDataVerifier) SendNonce(nonce string) {
	s.nonce <- nonce
}

func (s *signatureDataVerifier) SignatureData() signatureData {
	return <-s.signatureData
}

func (s *signatureDataVerifier) Verify() error {
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
	s.signatureData <- signatureData{
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
