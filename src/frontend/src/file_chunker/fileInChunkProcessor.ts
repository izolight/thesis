interface FileChunkDataCallback {
    (data: ArrayBuffer): void
}

interface ErrorCallback {
    (message: string): void
}

interface FileReaderOnLoadCallback {
    (event: ProgressEvent): void
}

interface ProcessingCompletedCallback {
    (): void
}

class Validate {
    public static notNull<T>(obj: T): obj is Exclude<T, null> {
        if (obj == null) {
            throw new ReferenceError(`Error: Object ${obj} was null`);
        }
        return true;
    }

    public static notUndefined<T>(obj: T): obj is Exclude<T, undefined> {
        if (obj === undefined) {
            throw new ReferenceError(`Error: Object ${obj} was undefined`);
        }
        return true;
    }

    public static notNullNotUndefined<T>(obj: T): obj is NonNullable<T> {
        return Validate.notNull(obj) && Validate.notUndefined(obj);
    }
}

 // TODO use SubtleCrypto for small files, wasm fuck for large files

class FileInChunksProcessor {
    public readonly CHUNK_SIZE_IN_BYTES: number = 1024*1000*20;
    private readonly fileReader: FileReader;
    private readonly dataCallback: FileChunkDataCallback;
    private readonly errorCallback: ErrorCallback;
    private readonly processingCompletedCallback: ProcessingCompletedCallback;
    private buffer: Uint8Array | null = null;
    private start: number = 0;
    private end: number = this.start + this.CHUNK_SIZE_IN_BYTES;
    private inputFile: File | null = null;

    constructor(dataCallback: FileChunkDataCallback,
                errorCallback: ErrorCallback,
                processingCompletedCallback: ProcessingCompletedCallback) {
        this.fileReader = new FileReader();
        this.fileReader.onload = this.getFileReadOnLoadHandler();
        this.dataCallback = dataCallback;
        this.errorCallback = errorCallback;
        this.processingCompletedCallback = processingCompletedCallback;
        console.log("startHash()");
        // @ts-ignore
        startHash();
    }

    public processChunks(inputFile: File) {
        this.inputFile = inputFile;
        this.read(this.start, this.end);
    }

    public getFileFromElement(elementId: string): File | undefined {
        const filesElement = document.getElementById(elementId) as HTMLInputElement;
        if (Validate.notNull(filesElement) &&
            Validate.notNull(filesElement.files)) {
            if (filesElement.files.length < 0) {
                this.errorCallback("Too few files specified. Please select one file.")
            } else if (filesElement.files.length > 1) {
                this.errorCallback("Too many files specified. Please select one file.")
            } else {
                return filesElement.files[0];
            }
        }
    }

    private getFileReadOnLoadHandler(): FileReaderOnLoadCallback {
        return () => {
            if (Validate.notNull(this.inputFile)) {
                this.buffer = new Uint8Array((this.fileReader.result as ArrayBuffer));
                this.dataCallback(this.buffer);

                this.start = this.end;
                if (this.end > this.inputFile.size) {
                    this.end = this.inputFile.size;
                } else {
                    this.end = this.start + this.CHUNK_SIZE_IN_BYTES;
                }
                if (this.end != this.inputFile.size) {
                    this.read(this.start, this.end);
                } else if(this.end == this.inputFile.size) {
                    this.processingCompletedCallback()
                }
            }
        }
    }

    private read(start: number, end: number) {
        if (Validate.notNull(this.inputFile)) {
            this.fileReader.readAsArrayBuffer(this.inputFile.slice(start, end));
        }
    }
}

function handleData(data: ArrayBuffer) {
    // const resultElement = document.getElementById("result");
    // if (Validate.notNull(resultElement)) {
    //     resultElement.innerHTML = `${resultElement.innerHTML} <p>Got ${data.byteLength} bytes: ${data.slice(0, 10).toString()}...</p>`;
    // }
    const str: string = `progressiveHash() len: ${data.byteLength} value: ${data.slice(0, 5).toString()}...`;
    console.log(str);
    // @ts-ignore
    progressiveHash(data);
}

function handleError(message: string) {
    const errorElement = document.getElementById("error");
    if (Validate.notNull(errorElement)) {
        errorElement.innerHTML = `${errorElement.innerHTML} <p>${message}</p>`;
    }
}

function processFileButtonHandler() {
    const processor = new FileInChunksProcessor(handleData, handleError, processingCompleted);
    const file = processor.getFileFromElement("file");
    if (Validate.notNullNotUndefined(file)) {
        console.log(`Started at ${new Date().toLocaleTimeString()}`)
        processor.processChunks(file);
    }
}

function processingCompleted() {
    // @ts-ignore
    const result = getHash();
    console.log(`getHash() -> ${result}`);
    const resultElement = document.getElementById("result");
    if (Validate.notNull(resultElement)) {
        resultElement.innerHTML = `<p>Got: ${result} </p>`;
        console.log(`Ended at ${new Date().toLocaleTimeString()}`)
    }

}

(function() {
    // @ts-ignore
    if (!WebAssembly.instantiateStreaming) {
        // @ts-ignore
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        }
    }
    // @ts-ignore
    const go = new Go();
    let mod, inst;
    // @ts-ignore
    WebAssembly.instantiateStreaming(fetch("../../test.wasm"), go.importObject).then(
        // @ts-ignore
        async result => {
            mod = result.module;
            inst = result.instance;
            await go.run(inst);
        }
    );
})();
