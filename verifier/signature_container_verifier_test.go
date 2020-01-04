package verifier

import (
	"golang.org/x/crypto/ocsp"
	"testing"
)

func Test_ocspStatusString(t *testing.T) {
	type args struct {
		status int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"good", args{ocsp.Good}, "Good",},
		{"revoked", args{ocsp.Revoked}, "Revoked",},
		{"unknown", args{ocsp.Unknown}, "Unknown",},
		{"bla", args{99}, "ServerFailed",},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ocspStatusString(tt.args.status); got != tt.want {
				t.Errorf("ocspStatusString() = %v, want %v", got, tt.want)
			}
		})
	}
}