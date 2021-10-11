package renderer

import (
	"github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
)

func Portal(target dom.Element) func(children ...interface{}) *VNode {
	return func(children ...interface{}) *VNode {
		v := Fragment(children...)
		v.IsRoot = true
		v.Node = target
		return v
	}
}
