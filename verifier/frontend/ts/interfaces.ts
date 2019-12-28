export interface PostHashesRequest {
    readonly hashes: Array<string>;
}

export interface Providers {
    [key: string]: string
}

export interface PostVerifyResponse {
    readonly valid: boolean;
    readonly error: string;
    readonly id_token: IDToken;
    readonly signature: Signature;
    readonly signing_cert: SigningCert;
    readonly timestamp: Timestamp;
}

export interface IDToken {
    readonly Issuer: string;
    readonly Audience: Array<string>;
    readonly Subject: string;
    readonly Expiry: string;
    readonly IssuedAt: string;
    readonly Nonce: string;
    readonly email: string;
    readonly email_verified: boolean;
    readonly cert_chain: Array<CertChain>;
}

export interface Signature {
    readonly salted_hashes: Array<string>;
    readonly hash_algorithm: string;
    readonly mac_key: string;
    readonly mac_algorithm: string;
    readonly signature_level: string;
}

export interface SigningCert {
    readonly signer: string;
    readonly signer_email: string;
    readonly cert_chain: Array<CertChain>;
}

export interface Timestamp {
    readonly SigningTime: string;
    readonly cert_chain: Array<CertChain>;
}

export interface CertChain {
    readonly issuer: string;
    readonly subject: string;
    readonly not_before: string;
    readonly not_after: string;
    readonly ocsp_status?: string;
    readonly ocsp_generation_time?: string;
}

export interface PoorPeoplePersistence {
    readonly postHashesResponse: PostVerifyResponse,
    readonly hashes: Array<string>
}

export interface FileChunkDataCallback {
    (data: Uint8Array): void
}

export interface Base64Callback {
    (data: string): void
}

export interface ErrorCallback {
    (message: string): void
}

export interface FileReaderOnLoadCallback {
    (event: ProgressEvent): void
}

export interface ProgressCallback {
    (percentCompleted: number): void
}

export interface ProcessingCompletedCallback {
    (): void
}

export interface Callable {
    (): void
}

export interface Continuable {
    (next: Callable): void
}
