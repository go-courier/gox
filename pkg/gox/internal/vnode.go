package internal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-courier/gox/pkg/dom"
)

type contextKeyVNode struct{}

func ContextWithVNode(ctx context.Context, vnode *VNode) context.Context {
	return context.WithValue(ctx, contextKeyVNode{}, vnode)
}

func VNodeFromContext(ctx context.Context) *VNode {
	if vnode, ok := ctx.Value(contextKeyVNode{}).(*VNode); ok {
		return vnode
	}
	return nil
}

type Key string

type Ref struct {
	Current interface{}
}

func JSX(tpe Component, children ...interface{}) *VNode {
	v := &VNode{
		Type:  tpe,
		Attrs: Attrs{},
	}

	for i := range children {
		switch x := children[i].(type) {
		case Key:
			// must bind first
			v.Key = x
		case *VNode:
			// avoiding be treated a fmt.Stringer
		case func(children ...interface{}) *VNode:
			// auto call
			children[i] = x()
		case Component:
			// auto render
			children[i] = JSX(x)
		case []interface{}:
			children[i] = JSX(Fragment{}, x...)
		case string:
			children[i] = JSX(Text(x))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			children[i] = JSX(Text(fmt.Sprintf("%v", x)))
		case fmt.Stringer:
			children[i] = JSX(Text(x.String()))
		}
	}

	v.InputChildren = children

	return v
}

type VNode struct {
	Type  Component
	Key   Key
	Attrs Attrs
	Ref   *Ref

	InputChildren []interface{}
	Children      []interface{}

	Parent *VNode

	IsRoot bool
	Node   dom.Element
	update func(vn *VNode)
	hooks
}

func (v *VNode) String() string {
	return fmt.Sprintf("%#v", v.Type)
}

func SameComponent(type1 Component, type2 Component) bool {
	if typeE1, ok := type1.(Element); ok {
		if typeE2, ok := type2.(Element); ok {
			return typeE1 == typeE2
		}
	}

	t1 := reflect.TypeOf(type1)
	for t1.Kind() == reflect.Ptr {
		t1 = t1.Elem()
	}

	t2 := reflect.TypeOf(type2)
	for t2.Kind() == reflect.Ptr {
		t2 = t2.Elem()
	}

	return t1 == t2
}

func (v *VNode) MountedNode() dom.Element {
	mounted := v.Node
	if mounted == nil && v.Parent != nil {
		return v.Parent.MountedNode()
	}
	return mounted
}

func (v *VNode) Use(hook Hook) Hook {
	return v.hooks.use(hook)
}

func (v *VNode) OnUpdate(fn func(vn *VNode)) {
	v.update = fn
}

func (v *VNode) Update() {
	v.update(v)
}

func (v *VNode) WillRender(oldVNode *VNode) {
	if oldVNode != nil {
		v.hooks = oldVNode.hooks
	}
	v.hooks.hookInit()
}

func (v *VNode) WillMount(oldVNode *VNode) {
	if oldVNode != nil {
		v.IsRoot = oldVNode.IsRoot
		v.Node = oldVNode.Node
	}
}

func (v *VNode) DidMount() {
	if v.Ref != nil {
		if v.Node != nil {
			v.Ref.Current = v.Node
		} else {
			// in golang Type is also Component instance.
			v.Ref.Current = v.Type
		}
	}
	v.hooks.commit()
}

func (v *VNode) Destroy() error {
	if v.Ref != nil {
		v.Ref.Current = nil
	}
	v.Node = nil
	v.hooks.destroy()
	return nil
}
