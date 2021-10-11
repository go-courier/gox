//go:build js && wasm
// +build js,wasm

package dom

import (
	"github.com/go-courier/gox/pkg/jsgo"
)

type jsEvent struct {
	JSValue
}

func (e *jsEvent) Get(propName string) interface{} {
	return jsgo.ToGoValue(e.JSValue.Get(propName))
}

func (e *jsEvent) Set(propName string, v interface{}) {
	e.JSValue.Set(propName, v)
}

func (e *jsEvent) Target() EventTarget {
	return asElement(e.JSValue.Get("target"))
}

func (e *jsEvent) PreventDefault() {
	e.Call("preventDefault")
}

func (e *jsEvent) StopImmediatePropagation() {
	e.Call("stopImmediatePropagation")
}

func (e *jsEvent) StopPropagation() {
	e.Call("stopPropagation")
}
