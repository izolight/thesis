import {Validate} from "./validate";
import {PoorPeoplePersistence} from "./interfaces";
import {Http} from "./http";
import {q} from "./tsQuery";
import {TS} from "./fileInChunkProcessor";

interface SigningRequest {
    id_token: string,
    seed: string,
    salt: string,
    hashes: Array<string>
}

interface SigningResponse {
    signature: string
}

export class CB {
    public static handle() {
        const signingResponseString = localStorage.getItem('lolnogenerics');
        if (Validate.notNull(signingResponseString)) {
            const p: PoorPeoplePersistence = JSON.parse(signingResponseString);
            const req: SigningRequest = {
                id_token: this.getIdTokenFromGetParameters(),
                seed: p.postHashesResponse.seed,
                salt: p.postHashesResponse.salt,
                hashes: p.hashes
            };
            Http.request<SigningResponse>(
                "POST",
                "/api/v1/sign",
                JSON.stringify(req),
                (response) => CB.showDownloadButton(response),
                err => console.log("error " + err),
                "application/json"
            );
        }
    }

    public static getIdTokenFromGetParameters(): string {
        const paramList = window.location.href.split('#')[1].split('&');
        for (const s of paramList) {
            const param = s.split('=');
            if (param[0] == 'id_token') {
                return param[1];
            }
        }
        const msg = 'Found no id_token in url';
        console.trace(msg);
        throw new Error(msg);
    }

    public static showDownloadButton(response: SigningResponse) {
        const inputFilesArea = q("input-files-area");
        if (Validate.notNull(inputFilesArea)) {
            if (Validate.notNull(inputFilesArea.parentNode)) {
                const template = `<p class="lead">Signature creation successful!</p>
                         <a href="SIGNATUREURL" class="button btn btn-block btn-outline-primary">Download</a>`;
                const newButton = document.createElement('div');
                TS.killAllChildren(inputFilesArea);
                newButton.innerHTML = template
                    .replace('SIGNATUREURL', response.signature)
                inputFilesArea.parentNode.insertBefore(newButton, inputFilesArea);
            }
        }
    }

    public static showErrorMessage(error: string) {
        const inputFilesArea = q("input-files-area");
        if (Validate.notNull(inputFilesArea)) {
            if (Validate.notNull(inputFilesArea.parentNode)) {
                const template = `<p class="lead">Signature creation failed: ERROR</p>`;
                const newElement = document.createElement('div');
                TS.killAllChildren(inputFilesArea);
                newElement.innerHTML = template
                    .replace('ERROR', error);
                inputFilesArea.parentNode.insertBefore(newElement, inputFilesArea);
            }
        }
    }
}