function wasmProgressiveHash(input) {
    let str = "progressivehash with " + input.toString();
    console.log(str.slice(0, 100));
    progressiveHash(input);
}

function wasmStartHash() {
    console.log("startHash");
    startHash();
}

function wasmGetHash() {
    console.log("getHash");
    return getHash();
}

(function() {
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        }
    }
    const go = new Go();
    let mod, inst;
WebAssembly.instantiateStreaming(fetch("../../test.wasm"), go.importObject).then(
        async result => {
            mod = result.module;
            inst = result.instance;
            await go.run(inst);
        }
    );
})();
