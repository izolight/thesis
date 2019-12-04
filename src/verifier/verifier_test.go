package verifier_test

import (
	"github.com/golang/protobuf/proto"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"testing"
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
		{name: "first signatureFile", args: args{"signaturefile", ""}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := readFile(t, tt.args.file)
			signatureFile := &verifier.SignatureFile{}
			if err := proto.Unmarshal(file, signatureFile); err != nil {
				t.Fatalf("could not unmarshal signature to protobuf: %w", err)
			}

			if err := verifier.VerifySignatureFile(signatureFile, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("VerifySignatureFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}