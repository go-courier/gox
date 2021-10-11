//go:build js && wasm
// +build js,wasm

package httputil

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"net/textproto"
	"sync"
	"time"

	"github.com/go-courier/gox/pkg/browser"
	"github.com/go-courier/gox/pkg/browser/worker"
)

//go:embed xhr.worker.js
var xhrWorker []byte

var XHRWorkerURL = "data:application/javascript," + browser.EncodeURIComponent(string(xhrWorker))

func GetShortConnClientContext(ctx context.Context, timeout time.Duration, transports ...Transport) *http.Client {
	t := &jsHttpTransport{}

	client := &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	for i := range transports {
		client.Transport = transports[i](client.Transport)
	}

	return client
}

type jsHttpTransport struct {
}

func (jsHttpTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	headers := map[string][]interface{}{}

	for k, vv := range req.Header {
		headers[k] = make([]interface{}, len(vv))
		for i := range vv {
			headers[k][i] = vv[i]
		}
	}

	request := map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
		//"headers": headers,
	}

	if req.Body != nil {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		request["body"] = string(data)
	}

	w := worker.New(XHRWorkerURL, nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		data := <-w.Receiver()
		_ = w.Close()

		respData := data.(map[string]interface{})

		r := &http.Response{}

		r.StatusCode = int(respData["status"].(float64))
		h, _ := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(respData["header"].(string)))).ReadMIMEHeader()
		r.Header = http.Header(h)
		r.Body = io.NopCloser(bytes.NewBufferString(respData["body"].(string)))

		resp = r
		wg.Done()
	}()

	w.PostMessage(request)
	wg.Wait()

	return
}
