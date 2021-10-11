package renderer

import (
	"context"
	"fmt"

	"github.com/go-courier/gox/pkg/gox/internal"

	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
)

func CreateRoot(root Element) *Root {
	r := &Root{
		doc:  root.OwnerDocument(),
		root: Portal(root)(),
	}
	r.cq.Start()
	return r
}

type Root struct {
	cq   commitQueue
	doc  Doc
	root *VNode
}

func (r *Root) Close() error {
	return r.cq.Close()
}

func (r *Root) Render(ctx context.Context, vnode *VNode) error {
	nextRoot := Portal(r.root.Node)(vnode)
	r.patchVNode(ctx, r.root, nextRoot)
	r.root = nextRoot
	r.cq.ForceCommit()
	return nil
}

func (r *Root) Act(fn func()) {
	fn()
	r.cq.ForceCommit()
}

func (r *Root) sameVNode(vnode1 *VNode, vnode2 *VNode) bool {
	return internal.SameComponent(vnode1.Type, vnode2.Type) && vnode1.Key == vnode2.Key
}

func (r *Root) patchVNode(ctx context.Context, oldVNode *VNode, vnode *VNode) {
	if vnode == oldVNode {
		return
	}

	switch vnode.Type.(type) {
	case internal.Text:
		r.mount(ctx, oldVNode, vnode)
		r.cq.Dispatch(func() {
			vnode.DidMount()
		})
	case internal.Element, internal.Fragment:
		walkChildren(ctx, vnode, vnode.InputChildren...)
		r.mount(ctx, oldVNode, vnode)
		r.cq.Dispatch(func() {
			vnode.DidMount()
		})
	default:
		// only component need to render
		var doRender func(ctx context.Context, oldVNode *VNode, vnode *VNode)

		doRender = func(ctx context.Context, oldVNode *VNode, vnode *VNode) {
			childCtx := ctx

			if cp, ok := vnode.Type.(internal.ContextProvider); ok {
				childCtx = cp.GetChildContext(ctx)
			}

			vnode.WillRender(oldVNode)

			walkChildren(childCtx, vnode, internal.JSX(internal.Fragment{}, vnode.Type.Render(internal.ContextWithVNode(childCtx, vnode), vnode.InputChildren...)))

			r.mount(childCtx, oldVNode, vnode)

			rendered := *vnode

			vnode.OnUpdate(func(vn *VNode) {
				doRender(ctx, &rendered, vn)
			})

			r.cq.Dispatch(func() {
				vnode.DidMount()
			})
		}

		doRender(ctx, oldVNode, vnode)
	}
}

func (r *Root) mount(childCtx context.Context, oldVNode *VNode, vnode *VNode) {
	if oldVNode == nil {
		switch x := vnode.Type.(type) {
		case internal.Text:
			vnode.Node = r.createTextNode(string(x))
			return
		case internal.Element:
			vnode.Node = r.createElement(string(x))
			r.patchNodeAttrs(nil, vnode)
		}

		mounted := vnode.MountedNode()

		// TODO mv this after mount
		r.addVNodes(childCtx, mounted, nil, vnode.Children, 0, len(vnode.Children)-1)

		return
	}

	vnode.WillMount(oldVNode)

	switch x := vnode.Type.(type) {
	case internal.Text:
		if x != oldVNode.Type {
			r.setTextContent(vnode.Node, string(x))
		}
		return
	case internal.Element:
		r.patchNodeAttrs(oldVNode, vnode)
	}

	oldChildren := oldVNode.Children
	childVNodes := vnode.Children

	mounted := vnode.MountedNode()

	if len(oldChildren) != 0 && len(childVNodes) != 0 {
		r.patchVNodes(childCtx, mounted, oldChildren, childVNodes)
	} else if len(childVNodes) != 0 {
		r.addVNodes(childCtx, mounted, nil, childVNodes, 0, len(childVNodes)-1)
	} else if len(oldChildren) != 0 {
		r.removeVNodes(childCtx, mounted, oldChildren, 0, len(oldChildren)-1)
	}
}

func walkChildren(ctx context.Context, v *internal.VNode, children ...interface{}) {
	v.Children = children

	for i := range v.Children {
		switch x := v.Children[i].(type) {
		case internal.Key:
		case *internal.Ref:
			v.Ref = x
		case internal.Attrs:
			v.Attrs.Merge(x)
		case internal.CanAttrs:
			v.Attrs.Merge(x.Attrs(ctx))
		case *internal.VNode:
			x.Parent = v
		default:
			if x != nil {
				fmt.Printf("unsupported child %T in %T\n", x, v.Type)
			}
		}
	}
}

func indexVNodeList(l []interface{}, n, i int) *VNode {
	if maxIndex := n - 1; maxIndex >= i && i >= 0 {
		if vn, ok := l[i].(*VNode); ok {
			return vn
		}
	}
	return nil
}

func (r *Root) patchVNodes(ctx context.Context, parentNode Element, oldChildren []interface{}, newChildren []interface{}) {
	oldStartIdx := 0
	newStartIdx := 0
	oldN := len(oldChildren)
	oldEndIdx := oldN - 1
	oldStartVNode := indexVNodeList(oldChildren, oldN, 0)
	oldEndVNode := indexVNodeList(oldChildren, oldN, oldEndIdx)
	newN := len(newChildren)
	newEndIdx := newN - 1
	newStartVNode := indexVNodeList(newChildren, newN, 0)
	newEndVNode := indexVNodeList(newChildren, newN, newEndIdx)

	var oldKeyToIdx map[Key]int
	var idxInOld int
	var elmToMove *VNode
	var beforeNode Element

	for oldStartIdx <= oldEndIdx && newStartIdx <= newEndIdx {
		if oldStartVNode == nil {
			oldStartIdx++
			oldStartVNode = indexVNodeList(oldChildren, oldN, oldStartIdx) // VNode might have been moved left
		} else if oldEndVNode == nil {
			oldEndIdx--
			oldEndVNode = indexVNodeList(oldChildren, oldN, oldEndIdx)
		} else if newStartVNode == nil {
			newStartIdx++
			newStartVNode = indexVNodeList(newChildren, newN, newStartIdx)
		} else if newEndVNode == nil {
			newEndIdx--
			newEndVNode = indexVNodeList(newChildren, newN, newEndIdx)
		} else if r.sameVNode(oldStartVNode, newStartVNode) {
			r.patchVNode(ctx, oldStartVNode, newStartVNode)
			oldStartIdx++
			oldStartVNode = indexVNodeList(oldChildren, oldN, oldStartIdx)
			newStartIdx++
			newStartVNode = indexVNodeList(newChildren, newN, newStartIdx)
		} else if r.sameVNode(oldEndVNode, newEndVNode) {
			r.patchVNode(ctx, oldEndVNode, newEndVNode)
			oldEndIdx--
			oldEndVNode = indexVNodeList(oldChildren, oldN, oldEndIdx)
			newEndIdx--
			newEndVNode = indexVNodeList(newChildren, newN, newEndIdx)
		} else if r.sameVNode(oldStartVNode, newEndVNode) {
			// VNode moved right
			r.patchVNode(ctx, oldStartVNode, newEndVNode)
			r.insertBefore(parentNode, oldStartVNode.Node, oldEndVNode.Node.NextSibling())
			oldStartIdx++
			oldStartVNode = indexVNodeList(oldChildren, oldN, oldStartIdx)
			newEndIdx--
			newEndVNode = indexVNodeList(newChildren, newN, newEndIdx)
		} else if r.sameVNode(oldEndVNode, newStartVNode) {
			// VNode moved left
			r.patchVNode(ctx, oldEndVNode, newStartVNode)
			r.insertBefore(parentNode, oldEndVNode.Node, oldStartVNode.Node)

			oldEndIdx--
			oldEndVNode = indexVNodeList(oldChildren, oldN, oldEndIdx)
			newStartIdx++
			newStartVNode = indexVNodeList(newChildren, newN, newStartIdx)
		} else {
			if oldKeyToIdx == nil {
				oldKeyToIdx = createKeyToOldIdx(oldChildren, oldStartIdx, oldEndIdx)
			}

			idxInOld = oldKeyToIdx[newStartVNode.Key]

			if idxInOld == 0 {
				r.patchVNode(ctx, nil, newStartVNode)
				r.insertBefore(parentNode, newStartVNode.Node, oldStartVNode.Node)
			} else {
				elmToMove = indexVNodeList(oldChildren, oldN, idxInOld)
				if !internal.SameComponent(elmToMove.Type, newStartVNode.Type) {
					r.patchVNode(ctx, nil, newStartVNode)
					r.insertBefore(parentNode, newStartVNode.Node, oldStartVNode.Node)
				} else {
					r.patchVNode(ctx, elmToMove, newStartVNode)
					oldChildren[idxInOld] = nil
					r.insertBefore(parentNode, elmToMove.Node, oldStartVNode.Node)
				}
			}
			newStartIdx++
			newStartVNode = indexVNodeList(newChildren, newN, newStartIdx)
		}
	}

	if oldStartIdx <= oldEndIdx || newStartIdx <= newEndIdx {
		if oldStartIdx > oldEndIdx {
			beforeNode = nil
			if n := indexVNodeList(newChildren, newN, newEndIdx+1); n != nil {
				beforeNode = n.Node
			}
			r.addVNodes(ctx, parentNode, beforeNode, newChildren, newStartIdx, newEndIdx)
		} else {
			r.removeVNodes(ctx, parentNode, oldChildren, oldStartIdx, oldEndIdx)
		}
	}
}

func (r *Root) addVNodes(ctx context.Context, parentNode Element, beforeNode Element, vnodes []interface{}, startIdx int, endIdx int) {
	for startIdx <= endIdx {
		if vn, ok := vnodes[startIdx].(*VNode); ok {
			r.patchVNode(ctx, nil, vn)
			if !vn.IsRoot {
				r.insertBefore(parentNode, vn.Node, beforeNode)
			}
		}
		startIdx++
	}
}

func (r *Root) removeVNodes(ctx context.Context, parentNode Element, vnodes []interface{}, startIdx int, endIdx int) {
	for startIdx <= endIdx {
		if vn, ok := vnodes[startIdx].(*VNode); ok {
			if vn.Node == nil {
				r.removeVNodes(ctx, parentNode, vn.Children, 0, len(vn.Children)-1)
			} else {
				if !vn.IsRoot {
					r.removeChild(parentNode, vn.Node)
				}
			}
			_ = vn.Destroy()
		}
		startIdx++
	}
}

func (r *Root) patchNodeAttrs(oldVNode *VNode, vnode *VNode) {
	oldAttrs := Attrs{}
	attrs := Attrs{}

	if oldVNode != nil {
		if oldVNode.Attrs != nil {
			oldAttrs = oldVNode.Attrs
		}
	}

	if vnode != nil {
		if vnode.Attrs != nil {
			attrs = vnode.Attrs
		}
	}

	// update modified attributes, add new attributes
	for key := range attrs {
		cur := attrs[key]
		if old, ok := oldAttrs[key]; ok {
			if cur != old {
				r.setAttribute(vnode.Node, key, cur)
			}
		} else {
			r.setAttribute(vnode.Node, key, cur)
		}
	}

	for key := range oldAttrs {
		if _, ok := attrs[key]; !ok {
			r.removeAttribute(vnode.Node, key)
		}
	}
}

func createKeyToOldIdx(childVNodes []interface{}, beginIdx int, endIdx int) map[Key]int {
	keyToIndexMap := map[Key]int{}
	for i := beginIdx; i <= endIdx; i++ {
		if vn, ok := childVNodes[i].(*VNode); ok {
			key := vn.Key
			if key != "" {
				keyToIndexMap[key] = i
			}
		}
	}
	return keyToIndexMap
}

func (r *Root) insertBefore(parent, new, old Node) {
	r.cq.Dispatch(func() {
		parent.InsertBefore(new, old)
	})
}

func (r *Root) removeChild(parent, old Node) {
	r.cq.Dispatch(func() {
		parent.RemoveChild(old)
	})
}

func (r *Root) setAttribute(node Element, k string, v interface{}) {
	r.cq.Dispatch(func() {
		node.SetAttribute(k, v)
	})
}

func (r *Root) removeAttribute(node Element, key string) {
	r.cq.Dispatch(func() {
		node.RemoveAttribute(key)
	})
}

func (r *Root) setTextContent(n Node, text string) {
	r.cq.Dispatch(func() {
		n.SetTextContent(text)
	})
}

type elementDeffer struct {
	Element
}

func (e *elementDeffer) Underlying() interface{} {
	return e.Element
}

func (r *Root) createElement(tag string) Element {
	e := &elementDeffer{}
	r.cq.Dispatch(func() {
		e.Element = r.doc.CreateElement(tag)
	})
	return e
}

func (r *Root) createTextNode(text string) Element {
	e := &elementDeffer{}
	r.cq.Dispatch(func() {
		e.Element = r.doc.CreateTextNode(text)
	})
	return e
}
