rust-wasm-files = out/.gitignore out/hasher_bg.d.ts out/hasher_bg.wasm out/hasher.d.ts out/hasher.js out/package/json out/README.md
ts-files = out/fileInChunkProcessor.js out/validate.js out/webWorker.js

all: $(rust-wasm-files) $(ts-files)

clean:
	rm out/*

$(rust-wasm-files): rustwasm/hasher/src/*.rs
	wasm-pack build rustwasm/hasher/ -d ../../out

$(ts-files): src/file_chunker/fileInChunkProcessor.ts tsconfig.json
	tsc

run: all
	go run main.go
