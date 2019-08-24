function processFile() {
    const files = document.getElementById("file").files;
    if(files.length > 1) {
        pushError("Too many files selected. Please select one file.")
    }
    else if(files.length < 1) {
        pushError("Too few files selected. Please select one file.")
    }

    clearErrors();

    const file = files[0];
    console.log(file);
    const fileReader = new FileReader();
    const chunkSize = 1024 * 100;
    let start = 0;
    let end = start + chunkSize;

    fileReader.onload = function () {
        const buffer = new Uint8Array(fileReader.result);
        const resultElement = document.getElementById("result");
        resultElement.innerHTML = `${resultElement.innerHTML} <p>Read ${buffer.byteLength} bytes: ${buffer.slice(0, 10).toString()}....</p>`;
        start = end;
        if(end > file.size) {
            end = file.size;
        }
        else {
            end = start + chunkSize;
        }
        if(end !== file.size) {
            fileReader.readAsArrayBuffer(file.slice(start, end));
        }
    };

    fileReader.readAsArrayBuffer(file.slice(start, end));
}

function pushError(errorMessage) {
    const errorElement = document.getElementById("errors");
    errorElement.innerHTML = `${errorElement.innerHTML} <p>${errorMessage}</p>`;
}

function clearErrors() {
    document.getElementById("errors").innerHTML = "";
}
