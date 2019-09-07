package main

import (
  "crypto/sha256"
  "syscall/js"
)

var hasher Hasher

func consoleLog(this js.Value, in[]js.Value) interface{} {
	println(js.Global().
		Get("document").
		Call("getElementById", "test").
		Get("value").
		String())
	return this
}

func registerCallbacks() {
	js.Global().Set("consoleLog", js.FuncOf(consoleLog))
	js.Global().Set("progressiveHash", js.FuncOf(progressiveHash))
	js.Global().Set("startHash", js.FuncOf(startHash))
}

type Hasher struct {
  input chan []byte
  result chan []byte
}

func (h *Hasher) sha256Hash(){
  hash := sha256.New()
  for {
    data, more := <- h.input
    if more {
      hash.Write(data)
    } else {
      break
    }
  }
  h.result <- hash.Sum(nil)
  return
}

func progressiveHash(this js.Value, in []js.Value) interface{} {
  hasher.input <- []byte(in[0].String())
  return this
}

func startHash(this js.Value, in []js.Value) interface{} {
  hasher.result = make(chan []byte)
  hasher.input = make(chan []byte)
  return this
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	registerCallbacks()
	<-c
}
