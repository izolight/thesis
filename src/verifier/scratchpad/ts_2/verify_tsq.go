package main

import (
	"fmt"
	"go.mozilla.org/pkcs7"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide a file")
	}
	in := os.Args[1]
	data, err := ioutil.ReadFile(in)
	if err != nil {
		log.Fatalf("Could not read file %s: %s", in, err)
	}

	timestamp, err := pkcs7.ParseTSResponse(data)
	if err != nil {
		log.Fatalf("could not parse timestamp response: %s", err)
	}
	fmt.Println(timestamp.Time)
	fmt.Println(timestamp.SerialNumber)
	fmt.Printf("%x\n", timestamp.HashedMessage)
	fmt.Println(timestamp.HashAlgorithm)
	fmt.Println(timestamp.Accuracy)
	for _, c := range timestamp.Certificates {
		fmt.Printf("Subject: %s\n", c.Subject)
		fmt.Printf("Issuer: %s\n\n", c.Issuer)
	}
}