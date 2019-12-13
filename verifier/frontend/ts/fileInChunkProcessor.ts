import {Validate} from "./validate";
import {Sha256hasher} from "../pkg";
import {q} from "./tsQuery";
import {
    Callable,
    ErrorCallback,
    FileChunkDataCallback,
    FileReaderOnLoadCallback,
    PoorPeoplePersistence,
    PostHashesResponse,
    ProcessingCompletedCallback,
    ProgressCallback
} from "./interfaces";
import {Http} from "./http";
import {Queue} from "./callbackHellChainer";


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

    public static showSubmissionButton(hashList: Array<string>, base64File: string) {
        const inputFilesArea = q("input-files-area");
        if (Validate.notNull(inputFilesArea)) {
            inputFilesArea.innerHTML = `<p class="lead">Hashing completed. Continue when ready</p>
                         <button type="button" class="btn btn-block btn-outline-primary" id="submithashes">Submit for verifying</button>`;
            const btn = q("submithashes");
            if (Validate.notNull(btn)) {
                (btn as HTMLButtonElement).onclick = (_) => {
                    (btn as HTMLButtonElement).disabled = true;
                    this.submitHashes(hashList, base64File);
                }
            }
        }
    }

    public static submitHashes(hashList: Array<string>, base64File: string) {
            Http.request<PostHashesResponse>('POST',
                '/verify',
                JSON.stringify({
                    hash: hashList[0],
                    signature: base64File
                }),
                response => {
                    console.log(response);
                    const p: PoorPeoplePersistence = {
                        postHashesResponse: response,
                        hashes: hashList
                    };
                    localStorage.setItem('lolnogenerics', JSON.stringify(p));
                    this.showIdpLoginButtons(response);
                },
                err => console.log(`error ${err}`),
                'application/json');
            console.log(`POST ${hashList}`);
    }

    public static showIdpLoginButtons(response: PostHashesResponse) {
        const inputFilesArea = q("input-files-area");
        console.log('fu1');
        if (Validate.notNull(inputFilesArea)) {
            if (Validate.notNull(inputFilesArea.parentNode)) {
                console.log('fu2');
                const template = `<p class="lead">Please select whom to authenticate with</p>
                         <a href="IDPURL" class="button btn btn-block btn-outline-primary">IDPNAME</a>`;
                this.killAllChildren(inputFilesArea);
                for (const key in response.providers) {
                    console.log('fu3');
                    const newIdpButton = document.createElement('div');
                    newIdpButton.innerHTML = template
                        .replace('IDPURL', response.providers[key])
                        .replace('IDPNAME', key);
                    inputFilesArea.parentNode.insertBefore(newIdpButton, inputFilesArea);
                }
            }
        }
    }


    public static progressCallbackBuilder(file: File, index: number):
        ProgressCallback {
        const cardElement = q(`file.${index}`);
        if (Validate.notNull(cardElement)) {
            return (percentCompleted => {
                cardElement.innerHTML = this.renderCardTemplate(file, `${percentCompleted}%`);
            });
        } else {
            return (_) => {
                console.log(`cardElement file.${index} was null, cannot update progress`);
            }
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
                cardElement.innerHTML = this.renderCardTemplate(file, hash);
                next();
            }
        } else {
            return () => {
                console.log(`cardElement file.${index} was null, cannot update progress`);
            }
        }
    }

    public static renderCardTemplate(file: File, hashValue: string): string {
        return `<div class="card mb-4 box-shadow">
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
        console.log(elementId);
        const filesElement = q(elementId) as HTMLInputElement;
        console.log(filesElement);
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

    if (Validate.notNullNotUndefined(sigFileList)) {
        const file = sigFileList[0];
        let fileReader = new FileReader();
        fileReader.onload = () => {
            if (Validate.notNull(file)) {
                fileReader.readAsBinaryString(file);
                let base64File = btoa(fileReader.result as string);

                const hashersQueue = new Queue(() => {
                    TS.showSubmissionButton(hashList, base64File)
                });

                if (Validate.notNullNotUndefined(fileList)) {
                    for (let i = 0; i < fileList.length; i++) {
                        const file = fileList[i];
                        const cardDeck = q('cardarea');
                        if (Validate.notNullNotUndefined(cardDeck)) {
                            const newCard = document.createElement('div');
                            newCard.id = `file.${i}`;
                            newCard.innerHTML = TS.renderCardTemplate(file, 'Queued');
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
                    }
                    hashersQueue.start();
                }
            }
        }
    }

}

