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
    (startTime: Date, endTime: Date, fileSize: number): void
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
    private startTime?: Date;
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
        this.startTime = new Date();
        this.numChunks = Math.round(this.inputFile.size / this.CHUNK_SIZE_IN_BYTES);
        this.read(this.start, this.end);
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
                    if (Validate.notNullNotUndefined(this.startTime)) {
                        this.processingCompletedCallback(this.startTime, new Date(), this.inputFile.size);
                    }
                    this.progressCallback(100);
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
    const fileList = getFilesFromElement("file");
    if (Validate.notNullNotUndefined(fileList)) {
        for (let i = 0; i < fileList.length; i++) {
            const file = fileList[i];
            const cardDeck = q('cardarea');
            if (Validate.notNullNotUndefined(cardDeck)) {
                const newCard = document.createElement('fubar');
                newCard.innerHTML = renderCardTemplate(file);
                if(Validate.notNull(cardDeck.parentNode)) {
                    cardDeck.parentNode.insertBefore(newCard, cardDeck);
                }
            }


            // new FileInChunksProcessor((data) => {
            //         wasmHasher.update(new Uint8Array((data)));
            //     },
            //     errorHandlingCallback,
            //     progressCallback,
            //     (startTime, endTime, sizeInBytes) => {
            //         const hashStr = wasmHasher.hex_digest();
            //         wasmHasher.free();
            //         // TODO add hash here
            //     }
            // ).processChunks(fileList[i]);
        }
    }
}

function renderCardTemplate(file: File): string {
    return `<div class="card mb-4 box-shadow">
            <div class="card-header">
                <h4 class="my-0 font-weight-normal">FILENAME</h4>
            </div>
            <div class="card-body">
                <ul class="list-unstyled mt-3 mb-4">
                    <li>Size: FILESIZE KB</li>
                    <li>Type: FILETYPE</li>
                </ul>
                <button type="button" class="btn btn-block btn-danger">Remove file</button>
            </div>
        </div>`
        .replace("FILENAME", file.name)
        .replace("FILESIZE", (file.size / 1024).toString())
        .replace("FILETYPE", file.type != "" ? file.type : "Unknown");
}

function dateObjectToTimeString(date: Date): string {
    return date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds() + "." + date.getMilliseconds();
}

function getFilesFromElement(elementId: string): FileList | undefined {
    const filesElement = q(elementId) as HTMLInputElement;

    if (Validate.notNull(filesElement.files)) {
        if (filesElement.files.length < 0) {
            alert("Too few files selected. Please select at least one file.")
        } else {
            return filesElement.files;
        }
    }
}
