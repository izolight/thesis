package _go

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

var hasher hash.Hash

//go:export progressiveHash
func progressiveHash(in []byte) {
	fmt.Printf("Writing %d bytes to hash\n")
	hasher.Write(in)
}

//go:export startHash
func startHash() {
	hasher = sha256.New()
}

//go:export getHash
func getHash() string {
	h := hasher.Sum(nil)
	hashStr := fmt.Sprintf("%x", h)
	fmt.Printf("Hash: %s\n", hashStr)
	return hashStr
}

func waitForever() {
	c := make(chan struct{}, 0)
	<-c
}

func main() {
	fmt.Println("WASM Go Initialized")
	waitForever()
}
