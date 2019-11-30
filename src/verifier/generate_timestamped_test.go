package verifier

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"testing"
)

func TestGenerateTimestamped(t *testing.T) {
	data := readFile(t, "hello_world_response.tsr")
	intermediateCA := parsePEM(t, "SwissSign TSA Platinum CA 2017 - G22.pem")
	intermediateCAOCSPFile := readFile(t, "SwissSign TSA Platinum CA 2017 - G22.pem.ocsp")
	tsaCA := parsePEM(t, "SwissSign ZertES TSA UNIT CH-2018.pem")
	tsaCAOCSPFile := readFile(t, "SwissSign ZertES TSA UNIT CH-2018.pem.ocsp")

	ts := &Timestamped{
		Rfc3161Timestamp: data,
		LtvTimestamp: map[string]*LTV{
			fmt.Sprintf("%x", sha256.Sum256(intermediateCA.Raw)): {
				Ocsp: intermediateCAOCSPFile,
			},
			fmt.Sprintf("%x", sha256.Sum256(tsaCA.Raw)): {
				Ocsp: tsaCAOCSPFile,
			},
		},
	}

	msg, err := proto.Marshal(ts)
	if err != nil {
		t.Fatalf("could not marshal timestmaped: %s", err)
	}

	if err = ioutil.WriteFile("testdata/hello_world_response.tsr.data", msg, 0644); err != nil {
		t.Fatalf("could not write file: %s", err)
	}
}
