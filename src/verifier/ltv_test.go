package verifier

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestVerifyLTV(t *testing.T) {
	rootCA := parsePEM(t, "SwissSign Platinum CA - G2.pem")
	intermediateCA := parsePEM(t, "SwissSign TSA Platinum CA 2017 - G22.pem")

	intermediateCAOCSPFile := readFile(t, "SwissSign TSA Platinum CA 2017 - G22.pem.ocsp")
	tsaCAOCSPFile := readFile(t, "SwissSign ZertES TSA UNIT CH-2018.pem.ocsp")

	tests := []struct{
		name string
		cert *x509.Certificate
		ltv map[string]*LTV
		wantErr bool
	}{
		{
			name: "root CA",
			cert: rootCA,
			ltv: nil,
			wantErr:false,
		},
		{
			name: "intermediate CA without ltv info",
			cert: intermediateCA,
			ltv: nil,
			wantErr: true,
		},
		{
			name: "intermediate CA with ocsp response",
			cert: intermediateCA,
			ltv: map[string]*LTV{
				fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): &LTV{
					Ocsp:                 intermediateCAOCSPFile,
				},
			},
			wantErr: false,
		},
		{
			name: "intermediate CA with wrong ocsp resonse",
			cert: intermediateCA,
			ltv: map[string]*LTV{
				fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): &LTV{
					Ocsp:                 tsaCAOCSPFile,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyLTV(tt.cert, tt.ltv); err != nil != tt.wantErr {
				t.Errorf("VerifyLTV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func parsePEM(t *testing.T, filename string) *x509.Certificate {
	t.Helper()
	file := readFile(t, filename)
	filePEM, _ := pem.Decode(file)
	cert, err := x509.ParseCertificate(filePEM.Bytes)
	if err != nil {
		t.Errorf("could not parse pem: %s", err)
	}
	return cert
}