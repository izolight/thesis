package verifier_test

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"io/ioutil"
	"testing"
)

func TestVerifyTimestamp(t *testing.T) {
	intermediateCA := parsePEM(t, "SwissSign TSA Platinum CA 2017 - G22.pem")
	intermediateCAOCSPFile := readFile(t, "SwissSign TSA Platinum CA 2017 - G22.pem.ocsp")
	tsaCA := parsePEM(t, "SwissSign ZertES TSA UNIT CH-2018.pem")
	tsaCAOCSPFile := readFile(t, "SwissSign ZertES TSA UNIT CH-2018.pem.ocsp")

	ltvData := map[string]*verifier.LTV{
		fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
			Ocsp: intermediateCAOCSPFile,
		},
		fmt.Sprintf("%x", sha256.Sum256(tsaCA.Raw)): {
			Ocsp: tsaCAOCSPFile,
		},
	}

	type args struct {
		data []byte
		timestamps [][]byte
		verifyLTV bool
		ltvData map[string]*verifier.LTV
	}
	tests := []struct {
		name       string
		args args
		wantErr    bool
		expectedErr error
	}{
		{
			name: "valid single timestamp",
			args: args{
				data: []byte("hello world\n"),
				timestamps: [][]byte{
					readFile(t, "hello_world_response.tsr"),
				},
				verifyLTV: true,
				ltvData: ltvData,
			},

			wantErr: false,
		},
		{
			name: "nested timestamp",
			args: args{
				data: []byte("hello world\n"),
				timestamps: [][]byte{
					readFile(t, "hello_world_response.tsr"),
					readFile(t, "hello_world_response.tsr_response.tsr"),
				},
				verifyLTV: true,
				ltvData: ltvData,
			},
			wantErr: false,
		},
		{
			name: "missing inner timestamp",
			args:args{
				data:       []byte("hello world\n"),
				timestamps: [][]byte{
					readFile(t, "hello_world_response.tsr_response.tsr"),
				},
				verifyLTV:  false,
			},
			wantErr: true,
		},
		{
			name: "hash mismatch",
			args: args{
				data: []byte("hello world"),
				timestamps: [][]byte{
					readFile(t, "hello_world_response.tsr"),
				},
				verifyLTV:false,
			},
			wantErr: true,
			expectedErr: verifier.ErrHashMismatch,
		},
		{
			name:       "no timestamps",
			args: args{
				data:       []byte("hello world"),
				verifyLTV:  false,
			},
			wantErr:    true,
			expectedErr: verifier.ErrNoTimestamps,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := verifier.NewTimestampVerifier(tt.args.timestamps, tt.args.verifyLTV, tt.args.ltvData)
			verifier.SendData(tt.args.data)
			err := verifier.Verify()
			if err != nil != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !errors.Is(err, tt.expectedErr) && tt.expectedErr != nil {
				t.Errorf("expected %s, got %s", tt.expectedErr, err)
			}
		})
	}
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		t.Errorf("Could not read file: %s", err)
	}
	return data
}
