interface FileChunkDataCallback {
    (data: ArrayBuffer): void
}

interface ErrorCallback {
    (message: string): void
}

interface FileReaderOnLoadCallback {
    (event: ProgressEvent): void
}

class FileInChunksProcessor {
    public readonly CHUNK_SIZE_IN_BYTES: number = 1024*100;
    private readonly fileReader: FileReader;
    private readonly dataCallback: FileChunkDataCallback;
    private readonly errorCallback: ErrorCallback;
    private buffer: Uint8Array | null = null;
    private start: number = 0;
    private end: number = this.start + this.CHUNK_SIZE_IN_BYTES;
    private inputFile: File | null = null;

    constructor(dataCallback: FileChunkDataCallback, errorCallback: ErrorCallback) {
        this.fileReader = new FileReader();
        this.fileReader.onload = this.getFileReadOnLoadHandler();
        this.dataCallback = dataCallback;
        this.errorCallback = errorCallback;
    }

    public processChunks(inputFile: File) {
        this.inputFile = inputFile;
        this.read(this.start, this.end);
    }

    public getFileFromElement(elementId: string): File | undefined {
        const filesElement = document.getElementById(elementId);
        if(filesElement == null) {
            console.log(`Error: element with id ${elementId} not found`);
        }
        else {
            const htmlInputElement = filesElement as HTMLInputElement;
            if(htmlInputElement == null || htmlInputElement.files == null) {
                console.log(`Error: element ${filesElement} is not a HTMLInputElement`);
            }
            else {
                if(htmlInputElement.files.length < 0) {
                    this.errorCallback("Too few files specified. Please select one file.")
                }
                else if (htmlInputElement.files.length > 1) {
                    this.errorCallback("Too many files specified. Please select one file.")
                }
                else {
                    return htmlInputElement.files[0];
                }
            }
        }
    }

    private getFileReadOnLoadHandler(): FileReaderOnLoadCallback {
        return () => {
            if (this.inputFile == null) {
                console.log("Error: Input File was null");
            }
            else {
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
                }
            }
        }
    }

    private read(start: number, end: number) {
        if(this.inputFile == null) {
            console.log("Error: Input File was null");
        }
        else {
            this.fileReader.readAsArrayBuffer(this.inputFile.slice(start, end));
        }
    }
}

function handleData(data: ArrayBuffer) {
    const resultElement = document.getElementById("result");
    if(resultElement == null) {
        console.log("Result Element was null");
    }
    else {
        resultElement.innerHTML = `${resultElement.innerHTML} <p>Got ${data.byteLength} bytes: ${data.slice(0, 10).toString()}...</p>`;
    }
}

function handleError(message: string) {
    const errorElement = document.getElementById("error");
    if(errorElement == null) {
        console.log("Error Element was null");
    }
    else {
        errorElement.innerHTML = `${errorElement.innerHTML} <p>${message}</p>`;
    }
}

function processFileButtonHandler() {
    const processor = new FileInChunksProcessor(handleData, handleError);
    const file = processor.getFileFromElement("file");
    if(file == null) {
        console.log("Input file was null");
    }
    else {
        processor.processChunks(file);

    }
}
