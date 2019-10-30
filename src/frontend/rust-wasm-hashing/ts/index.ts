import {Validate} from "./validate";
import {processFileButtonHandler} from "./fileInChunkProcessor"

function start(Sha256hasher: typeof import('../pkg')) {
    const button = document.getElementById("process-file");
    if (Validate.notNull(button)) {
        button.onclick = function () {
            processFileButtonHandler(new Sha256hasher.Sha256hasher());
        }
    }
}

async function load() {
    start(await import("../pkg"));
}

load();
