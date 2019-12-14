import {Validate} from "./validate";
import {processFileButtonHandler} from "./fileInChunkProcessor"

function start(Sha256hasher: typeof import('../pkg')) {
    const processFileButton = document.getElementById("process-file");
    if (Validate.notNull(processFileButton)) {
        processFileButton.onclick = function () {
            (processFileButton as HTMLButtonElement).disabled = true;
            processFileButtonHandler(new Sha256hasher.Sha256hasher());
        }
    }
}

async function load() {
        start(await import("../pkg"));
}

load();
