//go:build js && wasm
// +build js,wasm

package jsgo

import "syscall/js"

var (
	array  = js.Global().Get("Array")
	object = js.Global().Get("Object")
)

func ToGoValue(v js.Value) interface{} {
	switch v.Type() {
	case js.TypeString:
		return v.String()
	case js.TypeBoolean:
		return v.Bool()
	case js.TypeNumber:
		return v.Float()
	case js.TypeObject:
		if array.Call("isArray", v).Bool() {
			n := v.Length()
			slice := make([]interface{}, n)
			for i := 0; i < n; i++ {
				slice[i] = ToGoValue(v.Index(i))
			}
			return slice
		}

		m := map[string]interface{}{}
		keys := object.Call("keys", v)
		for i := 0; i < keys.Length(); i++ {
			k := keys.Index(i).String()
			m[k] = ToGoValue(v.Get(k))
		}
		return m
	case js.TypeNull:
		return nil
	default:
		return v
	}
}
