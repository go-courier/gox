package gox

import (
	"context"

	"github.com/go-courier/gox/pkg/gox/internal"
)

type VNode = internal.VNode
type Elem = internal.Element
type Key = internal.Key
type Ref = internal.Ref
type Component = internal.Component
type Attrs = internal.Attrs

func Attr(k string, v interface{}) Attrs {
	return Attrs{k: v}
}

func Fragment(children ...interface{}) *VNode {
	return H(internal.Fragment{})(children...)
}

func JSX(c Component, children ...interface{}) *VNode {
	return internal.JSX(c, children...)
}

func H(c Component) func(children ...interface{}) *VNode {
	return func(children ...interface{}) *VNode {
		return internal.JSX(c, children...)
	}
}

func Provider(c func(ctx context.Context) context.Context) func(children ...interface{}) *VNode {
	return func(children ...interface{}) *VNode {
		return internal.JSX(internal.Provider(c), children...)
	}
}

func UseEffect(ctx context.Context, setup func() func(), deps []interface{}) {
	internal.VNodeFromContext(ctx).Use(&internal.EffectHook{
		Setup: setup,
		Deps:  deps,
	})
}

func UseMemo(ctx context.Context, setup func() interface{}, deps []interface{}) interface{} {
	h := internal.VNodeFromContext(ctx).Use(&internal.MemoHook{
		Setup: setup,
		Deps:  deps,
	}).(*internal.MemoHook)

	return h.Memorised()
}

func UseRef(ctx context.Context, initialValue interface{}) *Ref {

	h := internal.VNodeFromContext(ctx).Use(&internal.RefHook{
		Ref: Ref{
			Current: initialValue,
		},
	}).(*internal.RefHook)

	return &h.Ref
}

func UseState(ctx context.Context, defaultState interface{}) (state interface{}, setState func(v interface{})) {
	vn := internal.VNodeFromContext(ctx)

	hook := vn.Use(&internal.StateHook{
		State:         defaultState,
		OnStateChange: vn.Update,
	}).(*internal.StateHook)

	return hook.State, hook.SetState
}
