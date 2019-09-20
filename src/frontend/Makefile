rust-wasm-files = out/.gitignore out/hasher_bg.d.ts out/hasher_bg.wasm out/hasher.d.ts out/hasher.js out/package.json out/README.md
ts-files = out/fileInChunkProcessor.js out/validate.js out/webWorker.js

all: $(rust-wasm-files) $(ts-files) webpack

webpack:
	npm run webpack -- --progress --display-optimization-bailout --display-error-details

clean:
	rm -r out/*

$(rust-wasm-files): rustwasm/hasher/src/lib.rs
	wasm-pack build rustwasm/hasher/ -d ../../out

$(ts-files): src/file_chunker/fileInChunkProcessor.ts tsconfig.json
	npm run tsc

run: all
	go run main.go
