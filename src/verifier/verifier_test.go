package verifier_test

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"

	"github.com/golang/protobuf/proto"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
)

func TestVerifySignatureFile(t *testing.T) {
	type args struct {
		file string
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "first signatureFile", args: args{"signaturefile", "06180c7ede6c6936334501f94ccfc5d0ff828e57a4d8f6dc03f049eaad5fb308"}, wantErr: false},
	}

	rootCA := parsePEM(t, "rootCA.pem")

	logger := log.New()
	logger.SetLevel(log.FatalLevel)
	cfg := verifier.Config{
		Issuer:   "https://keycloak.thesis.izolight.xyz/auth/realms/master",
		ClientId: "thesis",
		AdditionalCerts: []*x509.Certificate{
			rootCA,
		},
		Logger: log.NewEntry(logger),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := readFile(t, tt.args.file)
			signatureFile := &verifier.SignatureFile{}
			if err := proto.Unmarshal(file, signatureFile); err != nil {
				t.Fatalf("could not unmarshal signature file to protobuf: %w", err)
			}

			s := verifier.NewSignatureVerifier(cfg)
			resp, err := s.VerifySignatureFile(signatureFile, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifySignatureFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			body, err := json.MarshalIndent(resp, "", "\t")
			if err != nil {
				t.Error(body)
			}
			fmt.Printf("%s\n", body)
		})
	}
}

func TestGenerateFile(t *testing.T) {
	for i := 0; i < 10; i++ {
		b := make([]byte, 4)
		rand.Read(b)
		sf := &verifier.SignatureFile{
			SignatureDataInPkcs7: b,
			Rfc3161InPkcs7:       [][]byte{b},
		}
		b, err := proto.Marshal(sf)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%x\n", b)
	}
}

func TestParseSignatureFile(t *testing.T) {
	signatureFile := &verifier.SignatureFile{}
	if err := proto.Unmarshal(readFile(t, "signaturefile"), signatureFile); err != nil {
		t.Fatal(err)
	}
}
