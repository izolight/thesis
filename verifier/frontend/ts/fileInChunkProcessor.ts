import {Validate} from "./validate";
import {Sha256hasher} from "../pkg";
import {q} from "./tsQuery";
import {
    Base64Callback,
    Callable,
    ErrorCallback,
    FileChunkDataCallback,
    FileReaderOnLoadCallback,
    PoorPeoplePersistence,
    PostVerifyResponse,
    ProcessingCompletedCallback,
    ProgressCallback
} from "./interfaces";
import {Http} from "./http";
import {Queue} from "./callbackHellChainer";
import has = Reflect.has;


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
                processingCompletedCallback: Callable) {
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


export class TS {
    public static killAllChildren(e: HTMLElement) {
        while (e.children.length > 0) {
            e.removeChild(e.children[0]);
        }
    }

    public static showSubmissionButton(hashList: Array<string>, base64List: Array<string>) {
        const inputFilesArea = q("input-files-area");
        if (Validate.notNull(inputFilesArea)) {
            inputFilesArea.innerHTML = `<p class="lead">Hashing completed. Continue when ready</p>
                         <button type="button" class="btn btn-block btn-outline-primary" id="submithashes">Submit for verifying</button>`;
            const btn = q("submithashes");
            if (Validate.notNull(btn)) {
                (btn as HTMLButtonElement).onclick = (_) => {
                    (btn as HTMLButtonElement).disabled = true;
                    this.submitHashes(hashList, base64List);
                }
            }
        }
    }

    public static submitHashes(hashList: Array<string>, base64List: Array<string>) {
            Http.request<PostVerifyResponse>('POST',
                'verify',
                JSON.stringify({
                    hash: hashList[0],
                    signature: base64List[0],
                }),
                response => {
                    console.log(response);
                    const p: PoorPeoplePersistence = {
                        postHashesResponse: response,
                        hashes: hashList
                    };
                    localStorage.setItem('lolnogenerics', JSON.stringify(p));
                    this.showSignatureResult(response);
                },
                err => console.log(`error ${err}`),
                'application/json');
            console.log(`POST ${hashList}`);
    }

    public static showSignatureResult(response: PostVerifyResponse) {
        const cardArea = q("file.0");
        if (Validate.notNull(cardArea)) {
            const template = `<div class="card mb-4 box-shadow">
            <div class="card-header">
                <h5 class="my-0 font-weight-normal">Result</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled mt-3 mb-4">
                    <li>Signature Level: LEVEL</li>
                    <li>Signing Time: TIME</li>
                    <li>Signer E-Mail: EMAIL</li>
                    <li>Valid: VALID</li>
                </ul>
            </div>
        </div>`;
            cardArea.innerHTML = cardArea.innerHTML + template
                .replace('LEVEL', String(response.signature_level))
                .replace('TIME', response.signature_time)
                .replace('EMAIL', response.signer_email)
                .replace('VALID', String(response.valid))
        }
    }


    public static progressCallbackBuilder(file: File, index: number):
        ProgressCallback {
        const cardElement = q(`file.${index}`);
        if (Validate.notNull(cardElement)) {
            return (percentCompleted => {
                cardElement.innerHTML = this.renderCardTemplate(file, `${percentCompleted}%`, index);
            });
        } else {
            return (_) => {
                console.log(`cardElement file.${index} was null, cannot update progress`);
            }
        }
    }

    public static base64CompletedBuilder(next: Callable,
                                         base64List: Array<string>,
                                         file: File,
                                         index: number,
                                         base64er: Base64er
    ): Callable {
        const cardElement = q(`file.${index}`);
        if (Validate.notNull(cardElement)) {
            return () => {
                const base64File = base64er.get();
                base64List.push(base64File);
                cardElement.innerHTML = cardElement.innerHTML + this.renderBase64Template(file, index);
                next();
            }
        } else {
            return () => {}
        }
    }

    public static processingCompletedBuilder(next: Callable,
                                             hashList: Array<string>,
                                             file: File,
                                             index: number,
                                             wasmHasher: Sha256hasher
    ): Callable {
        const cardElement = q(`file.${index}`);
        if (Validate.notNull(cardElement)) {
            return () => {
                const hash = wasmHasher.hex_digest();
                hashList.push(hash);
                const innerElement = q(`filecard.${index}`);
                if(Validate.notNull(innerElement)) {
                    innerElement.innerHTML = "";
                }
                cardElement.innerHTML = cardElement.innerHTML + this.renderCardTemplate(file, hash, index);
                next();
            }
        } else {
            return () => {
                console.log(`cardElement file.${index} was null, cannot update progress`);
            }
        }
    }
    public static renderBase64Template(file: File, index: number): string {
        return `<div class="card mb-4 box-shadow" id="sigfilecard.${index}">
            <div class="card-header">
                <h5 class="my-0 font-weight-normal">FILENAME</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled mt-3 mb-4">
                    <li>Size: FILESIZE</li>
                </ul>
<!--                <button type="button" class="btn btn-block btn-danger">Remove file</button>-->
            </div>
        </div>`
            .replace("FILENAME", file.name)
            .replace("FILESIZE", file.size < 1024 * 1024 ? `${Math.round(file.size / 1024)} KB` : `${Math.round(file.size / 1024 / 1024)} MB`)
    }

    public static renderCardTemplate(file: File, hashValue: string, index: number): string {
        return `<div class="card mb-4 box-shadow"  id="filecard.${index}">
            <div class="card-header">
                <h5 class="my-0 font-weight-normal">FILENAME</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled mt-3 mb-4">
                    <li>Size: FILESIZE</li>
                    <li>Type: FILETYPE</li>
                    <li>Hash: FILEHASH</li>
                </ul>
<!--                <button type="button" class="btn btn-block btn-danger">Remove file</button>-->
            </div>
        </div>`
            .replace("FILENAME", file.name)
            .replace("FILESIZE", file.size < 1024 * 1024 ? `${Math.round(file.size / 1024)} KB` : `${Math.round(file.size / 1024 / 1024)} MB`)
            .replace("FILETYPE", file.type != "" ? file.type : "Unknown")
            .replace("FILEHASH", hashValue);
    }

    public static getFilesFromElement(elementId: string): FileList | undefined {
        const filesElement = q(elementId) as HTMLInputElement;
        if (Validate.notNull(filesElement.files)) {
            if (filesElement.files.length < 0) {
                alert("Too few files selected. Please select at least one file.")
            } else {
                this.updateFilesArea("Wait for hashing to finish")
            }
            return filesElement.files;
        }
    }

    public static updateFilesArea(message: string) {
        const inputFilesArea = q("input-files-area");
        if (Validate.notNull(inputFilesArea)) {
            //inputFilesArea.innerHTML = `<p class="lead">${message}</p>`;
        }
    }

    public static errorHandlingCallback(message: string) {
        const errorElement = q("error");
        errorElement.innerHTML = `${errorElement.innerHTML} <p>${message}</p>`;
    }
}


export function processFileButtonHandler(wasmHasher: Sha256hasher) {
    const fileList = TS.getFilesFromElement("file");
    const sigFileList = TS.getFilesFromElement("signature");
    const hashList = new Array<string>();
    const base64List = new Array<string>();

    const hashersQueue = new Queue(() => {
        TS.showSubmissionButton(hashList, base64List);
    });

    if (Validate.notNullNotUndefined(fileList) && Validate.notNullNotUndefined(sigFileList)) {
        if (fileList.length != sigFileList.length) {
            console.log(fileList, sigFileList);
            return
        }
        for (let i = 0; i < fileList.length; i++) {
            const file = fileList[i];
            const cardDeck = q('cardarea');
            if (Validate.notNullNotUndefined(cardDeck)) {
                const newCard = document.createElement('div');
                newCard.id = `file.${i}`;
                newCard.innerHTML = TS.renderCardTemplate(file, 'Queued', i);
                if (Validate.notNull(cardDeck.parentNode)) {
                    cardDeck.parentNode.insertBefore(newCard, cardDeck);
                }
            }

            hashersQueue.add(
                (next: Callable) => {
                    new FileInChunksProcessor((data) => {
                            wasmHasher.update(new Uint8Array((data)));
                        },
                        TS.errorHandlingCallback,
                        TS.progressCallbackBuilder(file, i),
                        TS.processingCompletedBuilder(next as Callable, hashList, file, i, wasmHasher)
                    ).processChunks(fileList[i]);
                }
            );

            const sigFile = sigFileList[i];
            let base64er = new Base64er();
            hashersQueue.add(
                (next: Callable) => {
                    new Base64Processor((data) => {
                            base64er.update(data);
                        },
                        TS.base64CompletedBuilder(next as Callable, base64List, sigFile, i, base64er)
                    ).process(sigFileList[i]);
                }
            );
        }

        hashersQueue.start();
    }
}

class Base64Processor {
    private readonly fileReader: FileReader;
    private readonly dataCallback: Base64Callback;
    private readonly processingCompletedCallback: ProcessingCompletedCallback;
    private inputFile: File | null = null;

    constructor(dataCallback: Base64Callback,
                processingCompletedCallback: Callable) {
        this.dataCallback = dataCallback;
        this.processingCompletedCallback = processingCompletedCallback;
        this.fileReader = new FileReader();
        this.fileReader.onload = this.getFileReadOnLoadHandler();
    }

    private getFileReadOnLoadHandler(): FileReaderOnLoadCallback {
        return () => {
            if (Validate.notNull(this.inputFile)) {
                this.dataCallback((this.fileReader.result as string));
                this.processingCompletedCallback();
            }
        }
    }

    public process(inputFile: File) {
        this.inputFile = inputFile;
        if (Validate.notNull(this.inputFile)) {
            this.fileReader.readAsBinaryString(this.inputFile);
        }
    }
}

class Base64er {
    private base64File: string = "";

    public update(input: string) {
        this.base64File = btoa(input);
    }

    public get(): string {
        return this.base64File;
    }
}