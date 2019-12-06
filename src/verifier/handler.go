package verifier

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"time"
)

type verifyRequest struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"` // base64 encoded protobuf file
}

type verifyResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
	SignerEmail string `json:"signer_email"`
	SignatureLevel SignatureLevel `json:"signature_level"`
	SignatureTime time.Time `json:"signature_time"`
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var in verifyRequest
	resp := verifyResponse{
		Valid: false,
	}
	if r.Body == nil {
		errorHandler(w, errors.New("no request body"), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(body, &in); err != nil {
		errorHandler(w, err, http.StatusBadRequest)
		return
	}
	signatureBytes, err := base64.StdEncoding.DecodeString(in.Signature)
	if err != nil {
		errorHandler(w, fmt.Errorf("could not decode signature: %w", err), http.StatusBadRequest)
		return
	}
	signatureFile := &SignatureFile{}
	if err := proto.Unmarshal(signatureBytes, signatureFile); err != nil {
		errorHandler(w, fmt.Errorf("could not unmarshal signature to protobuf: %w", err), http.StatusBadRequest)
		return
	}

	resp, err = VerifySignatureFile(signatureFile, in.Hash)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(resp)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

func errorHandler(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	resp := verifyResponse{
		Valid: false,
		Error: err.Error(),
	}
	out, _ := json.Marshal(resp)
	w.Write(out)
}
