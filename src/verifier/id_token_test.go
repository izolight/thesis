package verifier

import (
	"testing"
	"time"
)

func TestVerifyIDToken(t *testing.T) {
	idTokenFile := readFile(t, "idtoken_okta")
	idTokenManipulatedFile := readFile(t, "idtoken_okta_manipulated")
	tests := []struct{
		name string
		verifier Verifier
		wantErr bool
	}{
		{
			name: "valid id token (okta)",
			verifier: idTokenVerifier{
				token:idTokenFile,
				issuer: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7",
				keys: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7/v1/keys",
				nonce: "50adfcb0-5852-44fc-a8b8-8cd4216d7564",
				notAfter: time.Date(2019, 11, 18, 21, 00, 20, 0, time.Local),
			},
			wantErr: false,
		},
		{
			name: "valid , but expired id token (okta)",
			verifier: idTokenVerifier{
				token:idTokenFile,
				issuer: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7",
				keys: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7/v1/keys",
				nonce: "50adfcb0-5852-44fc-a8b8-8cd4216d7564",
				notAfter: time.Date(2019, 11, 18, 21, 03, 20, 0, time.Local),
			},
			wantErr: true,
		},
		{
			name: "wrong nonce (okta)",
			verifier: idTokenVerifier{
				token:idTokenFile,
				issuer: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7",
				keys: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7/v1/keys",
				nonce: "50adfcb0-5852-44fc-a8a8-8cd4216d7564",
				notAfter: time.Date(2019, 11, 18, 21, 02, 20, 0, time.Local),
			},
			wantErr:  true,
		},
		{
			name: "manipulated id token (okta)",
			verifier: idTokenVerifier{
				token: idTokenManipulatedFile,
				issuer: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7",
				keys: "https://micah.okta.com/oauth2/aus2yrcz7aMrmDAKZ1t7/v1/keys",
				nonce: "50adfcb0-5852-44fc-a8b8-8cd4216d7564",
				notAfter: time.Date(2019, 11, 18, 21, 02, 20, 0, time.Local),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.verifier.Verify(); err != nil != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

