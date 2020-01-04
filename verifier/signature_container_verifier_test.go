package verifier

import (
	"crypto/x509"
	"github.com/sirupsen/logrus"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier/pb"
	"golang.org/x/crypto/ocsp"
	"testing"
	"time"
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

func TestSignatureContainerVerifier_Verify(t *testing.T) {
	type fields struct {
		container       []byte
		data            chan pb.SignatureData
		signingCertData chan signingCertData
		signingTime     chan time.Time
		additionalCerts []*x509.Certificate
		cfg             *Config
	}
	type args struct {
		verifyLTV bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "invalid p7",
			fields:  fields{
				container:       []byte("test"),
				cfg: &Config{
					Logger:          logrus.NewEntry(logrus.New()),
				},
			},
			args:    args{false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SignatureContainerVerifier{
				container:       tt.fields.container,
				data:            tt.fields.data,
				signingCertData: tt.fields.signingCertData,
				signingTime:     tt.fields.signingTime,
				additionalCerts: tt.fields.additionalCerts,
				cfg:             tt.fields.cfg,
			}
			if err := s.Verify(tt.args.verifyLTV); (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}