package verifier_test

import (
	"crypto/sha256"
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
	//timestampedFile := readFile(t, "hello_world_response.tsr.data")

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
	}{
		{
			name: "valid single timestamp",
			args: args{
				data: []byte("hello world\n"),
				timestamps: [][]byte{
					readFile(t, "hello_world_response.tsr"),
				},
				verifyLTV: true,
				ltvData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Ocsp: intermediateCAOCSPFile,
					},
					fmt.Sprintf("%x", sha256.Sum256(tsaCA.Raw)): {
						Ocsp: tsaCAOCSPFile,
					},
				},
			},

			wantErr: false,
		},
		{
			name: "nested timestamp",
			args: args{
				data: []byte("hello world\n"),
				timestamps: [][]byte{
					readFile(t, "hello_world_response.tsr"),
					readFile(t, "hello_world_response.tsr.data_response.tsr"),
				},
				verifyLTV: true,
				ltvData: map[string]*verifier.LTV{
					fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
						Ocsp: intermediateCAOCSPFile,
					},
					fmt.Sprintf("%x", sha256.Sum256(tsaCA.Raw)): {
						Ocsp: tsaCAOCSPFile,
					},
				},
			},
			wantErr: false,
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
		},
		{
			name:       "no timestamps",
			args: args{
				data:       []byte("hello world"),
				verifyLTV:  false,
			},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := verifier.NewTimestampVerifier(tt.args.timestamps, tt.args.verifyLTV, tt.args.ltvData)
			verifier.SendData(tt.args.data)
			if err := verifier.Verify(); err != nil != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
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
