package verifier

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type VerifyRequest struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"` // base64 encoded protobuf file
}

type VerifyResponse struct {
	Valid          bool           `json:"valid"`
	Error          string         `json:"error,omitempty"`
	SignerEmail    string         `json:"signer_email"`
	SignatureLevel SignatureLevel `json:"signature_level"`
	SignatureTime  time.Time      `json:"signature_time"`
}

func NewDefaultCfg(caFile []byte) Config {
	cfg := Config{
		Issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
		ClientId: "thesis",
	}
	filePEM, _ := pem.Decode(caFile)
	rootCA, err := x509.ParseCertificate(filePEM.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	cfg.AdditionalCerts = []*x509.Certificate{
		rootCA,
	}
	return cfg
}

type VerifyService struct {
	cfg Config
}

func NewVerifyService(cfg Config) *VerifyService {
	return &VerifyService{cfg: cfg}
}

func (v *VerifyService) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	logger := newRequestLogger()
	localCfg := v.cfg
	localCfg.Logger = logger
	logger.Info("received verify request")

	w.Header().Set("Content-Type", "application/json")
	var in VerifyRequest
	resp := VerifyResponse{
		Valid: false,
	}
	if r.Body == nil {
		errorHandler(w, logger, errors.New("no request body"), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, logger, err, http.StatusInternalServerError)
		return
	}
	logger.Trace("read request body")
	if err = json.Unmarshal(body, &in); err != nil {
		errorHandler(w, logger, err, http.StatusBadRequest)
		return
	}
	logger.WithFields(log.Fields{
		"request_body": log.Fields{
			"hash": in.Hash,
		},
	}).Info("unmarshaled request body")

	signatureBytes, err := base64.StdEncoding.DecodeString(in.Signature)
	if err != nil {
		errorHandler(w, logger, fmt.Errorf("could not decode signature: %w", err), http.StatusBadRequest)
		return
	}
	logger.Trace("decoded signature")

	signatureFile := &SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		errorHandler(w, logger, fmt.Errorf("could not unmarshal signature to protobuf: %w", err), http.StatusBadRequest)
		return
	}
	logger.Info("unmarshaled signature file")

	s := NewSignatureVerifier(localCfg)
	resp, err = s.VerifySignatureFile(signatureFile, in.Hash)
	logger.Info("verified signature file")
	if err != nil {
		errorHandler(w, logger, err, http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(resp)
	if err != nil {
		errorHandler(w, logger, err, http.StatusInternalServerError)
		return
	}
	w.Write(out)
	logger.WithFields(log.Fields{
		"response_body": out,
	}).Trace("wrote response body")
}

func errorHandler(w http.ResponseWriter, logger *log.Entry, err error, code int) {
	logger.Error(err)
	w.WriteHeader(code)
	resp := VerifyResponse{
		Valid: false,
		Error: err.Error(),
	}
	out, _ := json.Marshal(resp)
	w.Write(out)
}

func newRequestLogger() *log.Entry {
	requestId := make([]byte, 16)
	rand.Read(requestId)
	return log.WithField("request_id", base64.StdEncoding.EncodeToString(requestId))
}
