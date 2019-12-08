package testdata

import (
	"encoding/pem"
	mozilla "go.mozilla.org/pkcs7"
	"io/ioutil"
	"os"
	"testing"
)

func TestPKCS7(t *testing.T) {
	tests := []struct {
		filename string
	}{
		{
			filename: "innerpkcs7",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			file := readFile(t, tt.filename)
			p7, err := mozilla.Parse(file)
			if err != nil {
				t.Errorf("could not parse with mozilla: %s", err)
			}
			b := &pem.Block{
				Type:    "CERTIFICATE",
				Bytes: p7.Certificates[2].Raw,
			}
			pem.Encode(os.Stdout, b)
		})
	}
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not read file: %s", err)
	}
	return data
}