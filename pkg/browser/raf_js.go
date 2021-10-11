//go:build js && wasm
// +build js,wasm

package browser

import "syscall/js"

func RequestAnimationFrame(fn func()) {
	var jsFunc js.Func
	jsFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		jsFunc.Release()
		return nil
	})
	js.Global().Call("requestAnimationFrame", jsFunc)
}
