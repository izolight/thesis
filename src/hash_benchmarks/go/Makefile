js-out = out/fileInChunkProcessor.js out/webWorker.js
js-src = src/ts/fileInChunkProcessor.ts src/ts/webWorker.ts
wasm-out = out/test.wasm
wasm-src = src/go/main.go

all: $(wasm-out) $(js-out)

$(wasm-out): $(wasm-src)
	GOOS=js GOARCH=wasm go build -o out/test.wasm ./src/go/main.go
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" lib/wasm_exec.js

$(js-out): $(js-src) tsconfig.json
	npm run tsc
