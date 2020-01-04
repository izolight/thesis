package api_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/api"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyHandler(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{name: "valid signature file", filename: "signaturefile_06180c.json", wantErr: false},
		{name: "no request body", filename: "empty.json", wantErr: true},
		{name: "invalid base64 file", filename: "invalid_base64.json", wantErr: true},
		{name: "invalid hash", filename: "invalid_hash.json", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestJSON := readFile(t, tt.filename)

			req, err := http.NewRequest("POST", "/verify", bytes.NewReader(requestJSON))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			caFile := readFile(t, "rootCA.pem")
			cfg := verifier.NewDefaultCfg(caFile)
			verifySvc := api.NewVerifyService(cfg)
			verifySvc.VerifyHandler(rr, req)

			body, _ := ioutil.ReadAll(rr.Result().Body)
			resp := &verifier.VerifyResponse{}

			if err := json.Unmarshal(body, resp); err != nil {
				t.Fatal(err)
			}
			if !resp.Valid != tt.wantErr {
				t.Errorf("validation failure: %s", resp.Error)
			}
		})
	}
}

func TestGenerateReqJSON(t *testing.T) {
	generateReqJSON(t, "signaturefile", "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308")
}

func generateReqJSON(t *testing.T, filename, hash string) {
	t.Helper()
	signatureFile := readFile(t, filename)
	sigbase64 := base64.StdEncoding.EncodeToString(signatureFile)
	verifyRequest := api.VerifyRequest{
		Hash:      hash,
		Signature: sigbase64,
	}
	requestJson, err := json.Marshal(verifyRequest)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("testdata/" + filename + "_" + hash[:6] + ".json", requestJson, 0644)
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		t.Errorf("Could not read file: %s", err)
	}
	return data
}
