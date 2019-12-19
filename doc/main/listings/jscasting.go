func progressiveHash(this js.Value, in []js.Value) interface{} {
  array := in[0]
  buf := make([]byte, array.Get("length").Int())
  js.CopyBytesToGo(buf, array)
  return this
}
