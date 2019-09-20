package main

import (
	"crypto/sha256"
	"fmt"
	"syscall/js"
)

var hasher = sha256.New()

func registerCallbacks() {
	js.Global().Set("progressiveHash", js.FuncOf(progressiveHash))
	js.Global().Set("startHash", js.FuncOf(startHash))
	js.Global().Set("getHash", js.FuncOf(getHash))
}

func progressiveHash(this js.Value, in []js.Value) interface{} {
	array := in[0]
	buf := make([]byte, array.Get("length").Int())
	n := js.CopyBytesToGo(buf, array)
	fmt.Printf("Copied %d bytes\n", n)
	hasher.Write(buf)
	return this
}

func startHash(this js.Value, in []js.Value) interface{} {
	hasher = sha256.New()
	return this
}

func getHash(this js.Value, in []js.Value) interface{} {
	hash := hasher.Sum(nil)
	hashStr := fmt.Sprintf("%x", hash)
	fmt.Printf("Hash: %s\n", hashStr)

	return js.ValueOf(hashStr)
}

func waitForever() {
	c := make(chan struct{}, 0)
	<-c
}

func main() {
	fmt.Println("WASM Go Initialized")
	registerCallbacks()
	waitForever()
}
