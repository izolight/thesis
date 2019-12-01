package verifier_test

import (
	"crypto/sha256"
	"fmt"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"testing"
	"time"
)

func TestVerifyIDToken(t *testing.T) {
	idTokenFile := readFile(t, "idtoken_keycloak")
	idTokenManipulatedFile := readFile(t, "idtoken_keycloak_manipulated")
	intermediateCAOCSPFile := readFile(t, "SwissSign TSA Platinum CA 2017 - G22.pem.ocsp")
	tsaCAOCSPFile := readFile(t, "SwissSign ZertES TSA UNIT CH-2018.pem.ocsp")
	intermediateCA := parsePEM(t, "SwissSign TSA Platinum CA 2017 - G22.pem")
	tsaCA := parsePEM(t, "SwissSign ZertES TSA UNIT CH-2018.pem")

	jwkFile := readFile(t, "jwk.json")
	ltv := map[string]*verifier.LTV{
		fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
			Ocsp: intermediateCAOCSPFile,
		},
		fmt.Sprintf("%x", sha256.Sum256(tsaCA.Raw)): {
			Ocsp: tsaCAOCSPFile,
		},
	}

	type args struct {
		token    []byte
		issuer   string
		nonce    string
		clientId string
		notAfter time.Time
		key      []byte
		ltv      map[string]*verifier.LTV
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid id token",
			args: args{
				token:    idTokenFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbe45",
				clientId: "thesis",
				notAfter: time.Unix(1575021202, 0),
				key:      jwkFile,
				ltv:      ltv,
			},
			wantErr: false,
		},
		{
			name: "valid , but expired id token (okta)",
			args: args{
				token:    idTokenFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbe45",
				clientId: "thesis",
				notAfter: time.Now(),
				key:      jwkFile,
				ltv:      ltv,
			},
			wantErr: true,
		},
		{
			name: "wrong nonce (okta)",
			args: args{
				token:    idTokenFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbea5",
				clientId: "thesis",
				notAfter: time.Unix(1575021202, 0),
				key:      jwkFile,
				ltv:      ltv,
			},
			wantErr: true,
		},
		{
			name: "manipulated id token (okta)",
			args: args{
				token:    idTokenManipulatedFile,
				issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
				nonce:    "5093b0fb5a68144fd3fddda5156f232e975c6eb857cba5b5fd9d64b7b31bbea5",
				clientId: "thesis",
				notAfter: time.Unix(1575021202, 0),
				key:      jwkFile,
				ltv:      ltv,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := verifier.NewIDTokenVerifier(nil, nil, time.Now())
			if err != nil {
				t.Errorf("NewIDTokenVerifier error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := v.Verify(); err != nil != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
