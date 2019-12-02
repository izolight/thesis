package verifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type verifyRequest struct {
	Hash      string `json:"hash"`
	Signature string `json:"signature"` // base64 encoded protobuf file
}

type verifyResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var in verifyRequest
	resp := verifyResponse{
		Valid: false,
	}
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = "No request body"
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	if err = json.Unmarshal(body, &in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	file, err := decodeSignatureFile(in)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Error = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	if err = verifySignatureFile(file, in.Hash); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	resp.Valid = true
	out, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp.Error = err.Error()
		out, _ := json.Marshal(resp)
		w.Write(out)
		return
	}
	w.Write(out)
}
