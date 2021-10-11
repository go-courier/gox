//go:build js && wasm
// +build js,wasm

package dom

import (
	"reflect"
	"strings"
	"sync"
	"syscall/js"

	"github.com/go-courier/gox/pkg/jsgo"
)

type jsElement struct {
	JSValue
	listeners sync.Map
}

func (e *jsElement) Get(propName string) interface{} {
	return jsgo.ToGoValue(e.JSValue.Get(propName))
}

func (e *jsElement) Set(propName string, v interface{}) {
	e.JSValue.Set(propName, v)
}

func (e *jsElement) NodeType() NodeType {
	return NodeType(e.Get("nodeType").(float64))
}

func (e *jsElement) NodeName() string {
	return strings.ToLower(e.Get("nodeName").(string))
}

func (e *jsElement) OwnerDocument() Doc {
	return Document
}

func (e *jsElement) TextContent() string {
	return e.Get("textContent").(string)
}

func (e *jsElement) SetTextContent(d string) {
	e.Set("textContent", d)
}

func (e *jsElement) ParentNode() Node {
	return asElement(e.JSValue.Get("parentNode"))
}

func (e *jsElement) FirstChild() Node {
	return asElement(e.JSValue.Get("firstChild"))
}

func (e *jsElement) LastChild() Node {
	return asElement(e.JSValue.Get("lastChild"))
}

func (e *jsElement) PreviousSibling() Node {
	return asElement(e.JSValue.Get("previousSibling"))
}

func (e *jsElement) NextSibling() Node {
	return asElement(e.JSValue.Get("nextSibling"))
}

func asElement(v js.Value) Element {
	if v.IsNull() {
		return nil
	}
	return &jsElement{
		JSValue: v,
	}
}

func (e *jsElement) AppendChild(c Node) {
	child := UnWrap(c)
	if child == nil {
		return
	}
	e.Call("appendChild", child.(*jsElement).JSValue)
}

func (e *jsElement) InsertBefore(new, old Node) Node {
	newChild := UnWrap(new)
	oldChild := UnWrap(old)

	if newChild == nil {
		return nil
	}
	if oldChild == nil {
		e.Call("insertBefore", newChild.(*jsElement).JSValue, nil)
	} else {
		e.Call("insertBefore", newChild.(*jsElement).JSValue, oldChild.(*jsElement).JSValue)
	}

	return newChild.(*jsElement)
}

func (e *jsElement) RemoveChild(c Node) Node {
	child := UnWrap(c)
	if child == nil {
		return nil
	}
	e.Call("removeChild", child.(*jsElement).JSValue)
	return child.(*jsElement)
}

func (e *jsElement) SetAttribute(k string, value interface{}) {
	e.Call("setAttribute", k, value)
}

func (e *jsElement) GetAttribute(k string) interface{} {
	v := e.Call("getAttribute", k)
	return jsgo.ToGoValue(v)
}

func (e *jsElement) RemoveAttribute(k string) {
	e.Call("removeAttribute", k)
}

func (e *jsElement) GetAttributeNames() []string {
	v := e.Call("getAttributeNames")
	names := make([]string, v.Length())
	for i := 0; i < v.Length(); i++ {
		names[i] = v.Index(i).String()
	}
	return names
}

func (e *jsElement) AddEventListener(eventName string, handle func(Event), args ...interface{}) {
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handle(&jsEvent{JSValue: args[0]})
		return nil
	})

	if len(args) == 1 {
		e.Call("addEventListener", eventName, fn, args[0])
	} else {
		e.Call("addEventListener", eventName, fn)
	}
	e.listeners.Store(reflect.TypeOf(handle), fn)
}

func (e *jsElement) RemoveEventListener(eventName string, handle func(Event), args ...interface{}) {
	k := reflect.TypeOf(handle)
	if fn, ok := e.listeners.Load(k); ok {
		f := fn.(js.Func)
		if len(args) == 1 {
			e.Call("removeEventListener", eventName, f, args[0])
		} else {
			e.Call("removeEventListener", eventName, f)
		}
		e.listeners.Delete(k)
		f.Release()
	}
}

func (e *jsElement) DispatchEvent(event Event) {
	// TODO
	e.Call("dispatchEvent", event)
}
