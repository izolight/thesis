package main

import (
	"syscall/js"
)

func consoleLog(this js.Value, in[] js.Value) interface{} {
	println(js.Global().
		Get("document").
		Call("getElementById", "test").
		Get("value").
		String())
	return this
}

func registerCallbacks() {
	js.Global().Set("consoleLog", js.FuncOf(consoleLog))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	registerCallbacks()
	<-c
}