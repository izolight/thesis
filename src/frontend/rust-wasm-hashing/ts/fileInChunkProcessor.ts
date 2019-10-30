import {Validate} from "./validate";
import {Sha256hasher} from "../pkg";
import {q} from "./tsQuery";

interface FileChunkDataCallback {
    (data: Uint8Array): void
}

interface ErrorCallback {
    (message: string): void
}

interface FileReaderOnLoadCallback {
    (event: ProgressEvent): void
}

interface ProgressCallback {
    (percentCompleted: number): void
}

interface ProcessingCompletedCallback {
    (): void
}


class FileInChunksProcessor {
    public readonly CHUNK_SIZE_IN_BYTES: number = 1024 * 1000 * 20;
    private readonly fileReader: FileReader;
    private readonly dataCallback: FileChunkDataCallback;
    private readonly errorCallback: ErrorCallback;
    private readonly processingCompletedCallback: ProcessingCompletedCallback;
    private readonly progressCallback: ProgressCallback;
    private start: number = 0;
    private end: number = this.start + this.CHUNK_SIZE_IN_BYTES;
    private numChunks: number = 0;
    private chunkCounter: number = 0;
    private inputFile: File | null = null;

    constructor(dataCallback: FileChunkDataCallback,
                errorCallback: ErrorCallback,
                progressCallback: ProgressCallback,
                processingCompletedCallback: ProcessingCompletedCallback) {
        this.fileReader = new FileReader();
        this.fileReader.onload = this.getFileReadOnLoadHandler();
        this.dataCallback = dataCallback;
        this.errorCallback = errorCallback;
        this.processingCompletedCallback = processingCompletedCallback;
        this.progressCallback = progressCallback;
    }

    public processChunks(inputFile: File) {
        this.inputFile = inputFile;
        this.numChunks = Math.round(this.inputFile.size / this.CHUNK_SIZE_IN_BYTES);
        this.read(this.start, this.end);
    }

    public getFileFromElement(elementId: string): File | undefined {
        const filesElement = q(elementId) as HTMLInputElement;

        if (Validate.notNull(filesElement.files)) {
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
                this.dataCallback(new Uint8Array((this.fileReader.result as ArrayBuffer)));

                this.start = this.end;
                if (this.end < this.inputFile.size) {
                    this.chunkCounter++;
                    this.end = this.start + this.CHUNK_SIZE_IN_BYTES;
                    this.progressCallback(Math.round((this.chunkCounter / this.numChunks) * 100));
                    this.read(this.start, this.end);
                } else {
                    this.processingCompletedCallback();
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

function errorHandlingCallback(message: string) {
    const errorElement = q("error");
    errorElement.innerHTML = `${errorElement.innerHTML} <p>${message}</p>`;
}

function progressCallback(percentCompleted: number) {
    if (percentCompleted == 100) {
        q("progress").innerText = `Completed!`
    } else {
        q("progress").innerText = `Hashing: ${percentCompleted}%`
    }
}

export function processFileButtonHandler(wasmHasher: Sha256hasher) {
    const startElement = q("start");
    const startTime = new Date();
    startElement.innerText = "started at " + dateObjectToTimeString(startTime);

    const processor = new FileInChunksProcessor((data) => {
            wasmHasher.update(new Uint8Array((data)));
        },
        errorHandlingCallback,
        progressCallback,
        () => {
            const hashStr = wasmHasher.hex_digest();
            wasmHasher.free();
            const endTime = new Date();
            const duration = (endTime.getTime() - startTime.getTime()) / 1000;
            q("result").innerHTML = `<p>Ended at ${dateObjectToTimeString(endTime)} <br>Got: ${hashStr} <br>${duration} seconds elapsed</p>`;
        }
    );
    const file = processor.getFileFromElement("file");
    if (Validate.notNullNotUndefined(file)) {
        processor.processChunks(file);
    }
}

function dateObjectToTimeString(date: Date): string {
    return date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds() + "." + date.getMilliseconds();
}
