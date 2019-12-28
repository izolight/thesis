package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"io/ioutil"
	"net/http"
)

type VerifyRequest struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"` // base64 encoded protobuf file
}

type VerifyService struct {
	cfg verifier.Config
}

func NewVerifyService(cfg verifier.Config) *VerifyService {
	return &VerifyService{cfg: cfg}
}

func (v *VerifyService) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	logger := newRequestLogger()
	localCfg := v.cfg
	localCfg.Logger = logger
	logger.Info("received verify request")

	w.Header().Set("Content-Type", "application/json")
	var in VerifyRequest
	resp := verifier.VerifyResponse{
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

	signatureFile := &verifier.SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		errorHandler(w, logger, fmt.Errorf("could not unmarshal signature to protobuf: %w", err), http.StatusBadRequest)
		return
	}
	logger.Info("unmarshaled signature file")

	s := verifier.NewSignatureVerifier(localCfg)
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
	resp := verifier.VerifyResponse{
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
