package verifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type verifyRequest struct {
	Hash string `json:"hash"`
	Signature string `json:"signature"`// base64 encoded protobuf file
}

type verifyResponse struct {
	Hash string `json:"hash"`
	Valid bool `json:"valid"`
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var in verifyRequest
	if r.Body == nil {
		http.Error(w, "No request body", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	valid, err := verifySignature(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := json.Marshal(verifyResponse{
		Hash: in.Hash,
		Valid: valid,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}