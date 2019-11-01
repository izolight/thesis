package main

import (
	"crypto/sha256"
	"fmt"
	"syscall/js"
	"time"
)

var hasher = sha256.New()
var start time.Time

func registerCallbacks() {
	js.Global().Set("progressiveHash", js.FuncOf(progressiveHash))
	js.Global().Set("startHash", js.FuncOf(startHash))
	js.Global().Set("getHash", js.FuncOf(getHash))
}

func progressiveHash(this js.Value, in []js.Value) interface{} {
	array := in[0]
	buf := make([]byte, array.Get("length").Int())
	js.CopyBytesToGo(buf, array)
	hasher.Write(buf)
	return this
}

func startHash(this js.Value, in []js.Value) interface{} {
	start = time.Now()
	fmt.Printf("Start: %s\n", start.Format(time.RFC3339Nano))
	hasher = sha256.New()
	return this
}

func getHash(this js.Value, in []js.Value) interface{} {
	hash := hasher.Sum(nil)
	hashStr := fmt.Sprintf("%x", hash)
	fmt.Printf("Hash: %s\n", hashStr)
	end := time.Now()
	fmt.Printf("End: %s\n",end.Format(time.RFC3339Nano))
	fmt.Printf("Took %s\n", end.Sub(start))

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
