package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"syscall/js"
)

var hasher hash.Hash

func main() {
	hasher = sha256.New()
}

//go:export progressiveHash
func progressiveHash(in []byte) {
	hasher.Write(in)
}

//go:export startHash
func startHash() {
	hasher = sha1.New()
}

//go:export getHash
func getHash() interface{} {
	h := hasher.Sum(nil)
	hashStr := fmt.Sprintf("%x", h)
	fmt.Printf("Hash: %s\n", hashStr)

	return js.ValueOf(hashStr)
}