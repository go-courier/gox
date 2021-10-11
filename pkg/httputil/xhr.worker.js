self.onmessage = (e) => {
    const req = e.data
    const xhr = new XMLHttpRequest();
    const onloadend = () => {
        self.postMessage({
            "status": xhr.status,
            "header": xhr.getAllResponseHeaders(),
            "body": xhr.response,
        });
    };
    if ('onloadend' in xhr) {
        xhr.onloadend = onloadend;
    } else {
        xhr.onreadystatechange = function handleLoad() {
            if (!xhr || xhr.readyState !== 4) {
                return;
            }
            if (xhr.status === 0 && !(xhr.responseURL && xhr.responseURL.indexOf('file:') === 0)) {
                return;
            }
            setTimeout(onloadend);
        };
    }
    xhr.open(req.method, req.url, true);
    if (req.headers) {
        for (const k in req.headers) {
            [].concat(req.headers[k]).forEach((header) => {
                xhr.setRequestHeader(k, header);
            })
        }
    }
    xhr.send(req.body ? req.body : undefined);
};
