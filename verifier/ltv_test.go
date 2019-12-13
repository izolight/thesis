package verifier_test

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"testing"
)

func TestVerifyLTV(t *testing.T) {
	rootCA := parsePEM(t, "SwissSign Platinum CA - G2.pem")
	intermediateCA := parsePEM(t, "SwissSign TSA Platinum CA 2017 - G22.pem")

	intermediateCAOCSPFile := readFile(t, "SwissSign TSA Platinum CA 2017 - G22.pem.ocsp")
	tsaCAOCSPFile := readFile(t, "SwissSign ZertES TSA UNIT CH-2018.pem.ocsp")

	silverCA := parsePEM(t, "SwissSign Silver CA - G2.pem")
	revokedIntermediateCA := parsePEM(t, "SwissSign Personal Silver CA 2014 - G22.pem")
	revokedIntermediateOCSPFile := readFile(t, "SwissSign Personal Silver CA 2014 - G22.pem.ocsp")

	tests := []struct {
		name     string
		verifier verifier.LTVVerifier
		wantErr  bool
	}{
		{
			name: "root CA",
			verifier: verifier.LTVVerifier{
				Certs:   []*x509.Certificate{rootCA},
				LTVData: nil,
			},
			wantErr: false,
		},
		{
			name: "intermediate CA without verifyLTV info",
			verifier: verifier.LTVVerifier{
				Certs:   []*x509.Certificate{rootCA, intermediateCA},
				LTVData: nil,
			},
			wantErr: true,
		},
		{
			name: "intermediate CA with nil verifyLTV",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{rootCA, intermediateCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): nil,
				},
			},
			wantErr: true,
		},
		{
			name: "intermediate CA with nil ocsp",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{rootCA, intermediateCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Ocsp: nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "intermediate CA with crl",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{rootCA, intermediateCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Crl: []byte("test"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "intermediate CA with ocsp response",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{rootCA, intermediateCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Ocsp: intermediateCAOCSPFile,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "intermediate CA with ocsp response and different ca order",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{intermediateCA, rootCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Ocsp: intermediateCAOCSPFile,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "intermediate CA with wrong ocsp response",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{rootCA, intermediateCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Ocsp: tsaCAOCSPFile,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "revoked ca",
			verifier: verifier.LTVVerifier{
				Certs: []*x509.Certificate{silverCA, revokedIntermediateCA},
				LTVData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(revokedIntermediateCA.Raw)): {
						Ocsp: revokedIntermediateOCSPFile,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.verifier.Verify(); err != nil != tt.wantErr {
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
