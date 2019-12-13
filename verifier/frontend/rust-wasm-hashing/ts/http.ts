export class Http {
    public static request<T>(
        method: 'GET' | 'POST',
        url: string,
        content?: string,
        callback?: (response: T) => void,
        errorCallback?: (err: any) => void,
        contentType?: string
    ) {
        const request = new XMLHttpRequest();
        const myHost = window.location.href.replace(window.location.pathname, '').replace(window.location.hash, '');
        const targetUrl = `${myHost}/${url}`;
        console.log(method + " " + targetUrl);
        request.open(method, targetUrl, true);
        request.onload = function () {
            if (this.status >= 200 && this.status < 400) {
                callback && callback(JSON.parse(this.response) as T);
            } else {
                errorCallback && errorCallback(this.responseText);
            }
        };

        request.onerror = function (err) {
            errorCallback && errorCallback(err);
        };

        if (method === 'POST') {
            if (contentType != null) {
                request.setRequestHeader(
                    'Content-Type',
                    contentType
                );

            }
        }
        request.send(content);
    }

    public static requestPromise<T>(
        method: 'GET' | 'POST',
        url: string,
        content ?: string
    ): Promise<T> {
        return new Promise<T>((resolve, reject) => {
            this.request(method, url, content, resolve, reject);
        });
    }
}
