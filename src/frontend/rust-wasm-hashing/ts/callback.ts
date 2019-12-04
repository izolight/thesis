import {Validate} from "./validate";
import {PoorPeoplePersistence} from "./interfaces";
import {Http} from "./http";

interface SigningRequest {
    idtoken: string,
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
        if(Validate.notNull(signingResponseString)) {
            const p: PoorPeoplePersistence = JSON.parse(signingResponseString);
            const req: SigningRequest = {
                idtoken : this.getIdTokenFromGetParameters(),
                seed: p.postHashesResponse.seed,
                salt: p.postHashesResponse.salt,
                hashes: p.hashes
            };
            Http.request<SigningResponse>(
                "POST",
                "/api/v1/sign",
                JSON.stringify(req),
                (response) => {
                    console.log(response)
                },
                err => console.log("error " + err),
                "application/json"
            );
        }
    }

    public static getIdTokenFromGetParameters(): string {
        const paramList = window.location.href.split('#')[1].split('&');
        for(const s of paramList) {
            const param = s.split('=');
            if(param[0] == 'id_token') {
                return param[1];
            }
        }
        const msg = 'Found no id_token in url';
        console.trace(msg);
        throw new Error(msg);
    }
}