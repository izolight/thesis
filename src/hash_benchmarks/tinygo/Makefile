js-out = out/fileInChunkProcessor.js out/webWorker.js
js-src = src/ts/fileInChunkProcessor.ts src/ts/webWorker.ts
wasm-out = out/test.wasm
wasm-src = src/go/main.go

all: $(wasm-out) $(js-out)

$(wasm-out): $(wasm-src)
	tinygo build -o out/test.wasm -target wasm ./src/go/main.go

$(js-out): $(js-src) tsconfig.json
	npm run tsc
