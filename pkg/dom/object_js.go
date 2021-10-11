//go:build js && wasm
// +build js,wasm

package dom

import (
	"syscall/js"
)

type JSValue = js.Value
