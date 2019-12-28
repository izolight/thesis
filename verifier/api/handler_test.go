package api_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyHandler(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		hash     string
		wantErr  bool
	}{
		{name: "valid signature file", filename: "signaturefile", hash: "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestJSON := generateReqJSON(t, tt.filename, tt.hash)

			req, err := http.NewRequest("POST", "/verify", bytes.NewReader(requestJSON))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			caFile := readFile(t, "rootCA.pem")
			cfg := verifier.NewDefaultCfg(caFile)
			verifySvc := verifier.NewVerifyService(cfg)
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

func generateReqJSON(t *testing.T, filename, hash string) []byte {
	t.Helper()
	signatureFile := readFile(t, filename)
	sigbase64 := base64.StdEncoding.EncodeToString(signatureFile)
	verifyRequest := verifier.VerifyRequest{
		Hash:      hash,
		Signature: sigbase64,
	}
	requestJson, err := json.Marshal(verifyRequest)
	if err != nil {
		t.Fatal(err)
	}
	return requestJson
}
