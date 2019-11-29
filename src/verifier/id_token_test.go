package verifier

import (
	"testing"
	"time"
)

func TestVerifyIDToken(t *testing.T) {
	idTokenFile := readFile(t, "idtoken_keycloak")
	idTokenManipulatedFile := readFile(t, "idtoken_keycloak_manipulated")
	tests := []struct{
		name string
		verifier Verifier
		wantErr bool
	}{
		{
			name: "valid id token",
			verifier: idTokenVerifier{
				token:    idTokenFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				keys:     "https://keycloak.thesis.izolight.xyz/auth/realms/master/protocol/openid-connect/certs",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbe45",
				clientId: "thesis",
				notAfter: time.Unix(1575021202, 0),
			},
			wantErr: false,
		},
		{
			name: "valid , but expired id token (okta)",
			verifier: idTokenVerifier{
				token:    idTokenFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				keys:     "https://keycloak.thesis.izolight.xyz/auth/realms/master/protocol/openid-connect/certs",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbe45",
				clientId: "thesis",
				notAfter: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "wrong nonce (okta)",
			verifier: idTokenVerifier{
				token:    idTokenFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				keys:     "https://keycloak.thesis.izolight.xyz/auth/realms/master/protocol/openid-connect/certs",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbea5",
				clientId: "thesis",
				notAfter: time.Unix(1575021202, 0),
			},
			wantErr:  true,
		},
		{
			name: "manipulated id token (okta)",
			verifier: idTokenVerifier{
				token:    idTokenManipulatedFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				keys:     "https://keycloak.thesis.izolight.xyz/auth/realms/master/protocol/openid-connect/certs",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbea5",
				clientId: "thesis",
				notAfter: time.Unix(1575021202, 0),
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

