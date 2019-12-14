export interface PostHashesRequest {
    readonly hashes: Array<string>;
}

export interface Providers {
    [key: string]: string
}

export interface PostVerifyResponse {
    readonly signature_level: number;
    readonly signature_time: string;
    readonly signer_email: string;
    readonly valid: boolean;
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
