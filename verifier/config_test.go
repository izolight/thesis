package verifier_test

import (
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"testing"
)

func TestNewDefaultCfg(t *testing.T) {
	type args struct {
		caFile []byte
	}
	notAPem := readFile(t, "notA.pem")
	invalidPem := readFile(t, "invalid.pem")

	tests := []struct {
		name string
		args args
	}{
		{
			name: "nil ca file",
			args: args{caFile:nil},
		},
		{
			name: "not a pem",
			args: args{caFile: notAPem},
		},
		{
			name: "invalid pem",
			args: args{caFile: invalidPem},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier.NewDefaultCfg(tt.args.caFile)
		})
	}
}