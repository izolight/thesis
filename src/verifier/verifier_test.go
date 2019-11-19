package verifier

import (
	"reflect"
	"testing"
)

/*func Test_verifyHashes(t *testing.T) {
	type args struct {
		data *SignatureData
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyHashes(tt.args.data, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("verifyHashes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}*/

func Test_verifyIDToken(t *testing.T) {
	type args struct {
		data *SignatureData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyIDToken(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("verifyIDToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifySignature(t *testing.T) {
	type args struct {
		container *SignatureContainer
	}
	tests := []struct {
		name    string
		args    args
		want    *SignatureData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifySignature(tt.args.container)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifySignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("verifySignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verifySignatureFile(t *testing.T) {
	type args struct {
		in verifyRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifySignatureFile(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("verifySignatureFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}