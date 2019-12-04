import {Validate} from "./validate";
import {processFileButtonHandler} from "./fileInChunkProcessor"
import {CB} from "./callback";

function start(Sha256hasher: typeof import('../pkg')) {
    const processFileButton = document.getElementById("process-file");
    if (Validate.notNull(processFileButton)) {
        processFileButton.onclick = function () {
            processFileButtonHandler(new Sha256hasher.Sha256hasher());
        }
    }
}

async function load() {
    if(document.location.pathname.includes("index")) {
        start(await import("../pkg"));
    }
    else if(document.location.pathname.includes("callback")){
        CB.handle();
    }
}

load();
