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
	ValidText string `json:"valid_text"`
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var in verifyRequest
	resp := verifyResponse{
		Hash: in.Hash,
		Valid: false,
	}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.ValidText = "No request body"
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.ValidText = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	err = json.Unmarshal(body, &in)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.ValidText = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	err = verifySignatureFile(in)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.ValidText = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	resp.ValidText = "Ok"
	resp.Valid = true
	out, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.ValidText = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	w.Write(out)
}