import {Validate} from "./validate";

interface WebWorkerResponseCallback {
    (data: WebWorkerResponse): void // TODO type?
}

type WebWorkerMethod = "progressiveHash" | "startHash" | "getHash";

interface WebWorkerMethodCall {
    method: WebWorkerMethod;
    data: any | null;
}

interface WebWorkerResponse {
    data: Exclude<any, null>;
}

class WebWorker {
    private readonly worker: Worker;

    public static isWebWorkerSupportPresent(): boolean {
        return typeof(Worker) !== "undefined";
    }

    constructor(codeToRun: string,
                webWorkerResponseCallback: WebWorkerResponseCallback) {
        Validate.notNull(codeToRun);
        Validate.notNull(webWorkerResponseCallback);

        this.worker = new Worker(codeToRun);
        this.worker.onmessage = webWorkerResponseCallback;
    }

    public callMethod(methodCall: WebWorkerMethodCall): void {
        this.worker.postMessage(methodCall);
    }
}