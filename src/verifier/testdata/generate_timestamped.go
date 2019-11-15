package main

import (
	"github.com/golang/protobuf/proto"
	"gitlab.ti.bfh.ch/hirtp1/thesis/src/verifier"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Could not read file: %s", err)
	}

	ts := &verifier.Timestamped{
		Rfc3161Timestamp:     data,
		LtvTimestamp:         nil,
	}

	msg, err := proto.Marshal(ts)
	if err != nil {
		log.Fatalf("could not marshal timestmaped: %s", err)
	}

	err = ioutil.WriteFile(os.Args[1] + ".data", msg, 0644)
	if err != nil {
		log.Fatalf("could not write file: %s", err)
	}
}