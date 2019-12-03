import {Validate} from "./validate";
import {processFileButtonHandler} from "./fileInChunkProcessor"

function start(Sha256hasher: typeof import('../pkg')) {
    const processFileButton = document.getElementById("process-file");
    // TODO bind to submit button
    if (Validate.notNull(processFileButton)) {
        processFileButton.onclick = function () {
            processFileButtonHandler(new Sha256hasher.Sha256hasher());
        }
    }
}

async function load() {
    start(await import("../pkg"));
}

load();
