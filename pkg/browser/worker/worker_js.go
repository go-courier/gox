//go:build js && wasm
// +build js,wasm

package worker

import (
	"syscall/js"

	"github.com/go-courier/gox/pkg/jsgo"
)

var (
	worker = js.Global().Get("Worker")
)

func New(url string, options map[string]interface{}) Worker {
	w := &jsWorker{
		worker:   worker.New(url, options),
		receiver: make(chan interface{}),
	}

	w.callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.receiver <- jsgo.ToGoValue(args[0].Get("data"))
		w.callback.Release()
		return nil
	})

	w.worker.Call("addEventListener", "message", w.callback)

	return w
}

type jsWorker struct {
	worker   js.Value
	receiver chan interface{}
	callback js.Func
}

func (w *jsWorker) Receiver() <-chan interface{} {
	return w.receiver
}

func (w *jsWorker) PostMessage(v interface{}) {
	w.worker.Call("postMessage", v)
}

func (w *jsWorker) Close() error {
	w.worker.Call("terminate")
	w.callback.Release()
	close(w.receiver)
	return nil
}
