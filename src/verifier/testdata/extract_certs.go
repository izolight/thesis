package testdata

import (
	"bytes"
	"encoding/pem"
	"go.mozilla.org/pkcs7"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Could not read file: %s", err)
	}
	ts, err := pkcs7.ParseTSResponse(data)
	if err != nil {
		log.Fatalf("could not parse timestamp response: %s", err)
	}
	for _, cert := range ts.Certificates {
		block := &pem.Block{
			Bytes: cert.Raw,
			Type: "CERTIFICATE",
		}
		var buf bytes.Buffer
		if err := pem.Encode(&buf, block); err != nil {
			log.Fatalf("could not encode pem: %s", err)
		}
		filename := cert.Subject.CommonName + ".pem"
		if err := ioutil.WriteFile(filename, buf.Bytes(), 0644); err != nil {
			log.Fatalf("could not write file %s: %s", filename, err)
		}
	}
}