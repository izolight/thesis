package verifier

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestVerifyTimestamp(t *testing.T) {
	intermediateCA := parsePEM(t, "SwissSign TSA Platinum CA 2017 - G22.pem")
	intermediateCAOCSPFile := readFile(t, "SwissSign TSA Platinum CA 2017 - G22.pem.ocsp")
	tsaCA := parsePEM(t, "SwissSign ZertES TSA UNIT CH-2018.pem")
	tsaCAOCSPFile := readFile(t, "SwissSign ZertES TSA UNIT CH-2018.pem.ocsp")

	tests := []struct {
		name       string
		data       []byte
		timestamps []*Timestamped
		wantErr    bool
	}{
		{
			name: "valid single timestamp",

			data: []byte("hello world\n"),
			timestamps: []*Timestamped{
				{
					Rfc3161Timestamp: readFile(t, "hello_world_response.tsr"),
					LtvTimestamp: map[string]*LTV{
						fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
							Ocsp: intermediateCAOCSPFile,
						},
						fmt.Sprintf("%x", sha256.Sum256(tsaCA.Raw)): {
							Ocsp: tsaCAOCSPFile,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nested timestamp",
			data: []byte("hello world\n"),
			timestamps: []*Timestamped{
				{
					Rfc3161Timestamp: readFile(t, "hello_world_response.tsr"),
					LtvTimestamp:     map[string]*LTV{},
				},
				{
					Rfc3161Timestamp: readFile(t, "hello_world_response.tsr.data_response.tsr"),
					LtvTimestamp:     map[string]*LTV{},
				},
			},
			wantErr: false,
		},
		{
			name: "hash mismatch",
			data: []byte("hello world"),
			timestamps: []*Timestamped{
				{
					Rfc3161Timestamp: readFile(t, "hello_world_response.tsr"),
					LtvTimestamp:     map[string]*LTV{},
				},
			},
			wantErr: true,
		},
		{
			name:       "no timestamps",
			data:       []byte("hello world"),
			timestamps: []*Timestamped{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := NewTimestampVerifier(tt.timestamps)
			verifier.sendData(tt.data)
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
