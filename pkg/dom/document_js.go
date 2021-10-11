//go:build js && wasm
// +build js,wasm

package dom

import (
	"syscall/js"

	"github.com/go-courier/gox/pkg/jsgo"
)

var (
	Document = &jsDocument{JSValue: js.Global().Get("document")}
)

type jsDocument struct {
	JSValue
}

func (e *jsDocument) Get(propName string) interface{} {
	return jsgo.ToGoValue(e.JSValue.Get(propName))
}

func (e *jsDocument) Set(propName string, v interface{}) {
	e.JSValue.Set(propName, v)
}

func (d *jsDocument) QuerySelector(s string) Element {
	jsVal := d.Call("querySelector", s)
	if jsVal.IsNull() {
		return nil
	}
	return &jsElement{
		JSValue: jsVal,
	}
}

func (d *jsDocument) CreateElement(tagName string) Element {
	return &jsElement{JSValue: d.Call("createElement", tagName)}
}

func (d *jsDocument) CreateTextNode(text string) Element {
	return &jsElement{JSValue: d.Call("createTextNode", text)}
}
