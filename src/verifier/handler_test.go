package verifier_test

import (
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyHandler(t *testing.T) {
	tests := []struct {
		name string
		filename string
		hash string
		wantErr bool
	}{
		{name: "valid signature file", filename: "signatureFile", hash: "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest("POST", "/verify", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(verifier.VerifyHandler)
			handler.ServeHTTP(rr, req)
		})
	}
}