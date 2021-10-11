//go:build !js
// +build !js

package dom

import (
	"sync"
)

type element struct {
	parent, firstChild, lastChild, prevSibling, nextSibling *element

	tagName     string
	nodeType    NodeType
	textContent string
	attributes  map[string]interface{}

	object

	rw sync.RWMutex
}

func (*element) OwnerDocument() Doc {
	return Document
}

func (e *element) NodeType() NodeType {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e == nil {
		return ERROR_NODE
	}
	return e.nodeType
}

func (e *element) NodeName() string {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.nodeType == TEXT_NODE {
		return "#text"
	}
	return e.tagName
}

func (e *element) TextContent() string {
	e.rw.RLock()
	defer e.rw.RUnlock()

	return e.textContent
}

func (e *element) SetTextContent(d string) {
	e.rw.Lock()
	defer e.rw.Unlock()

	e.textContent = d
}

func (e *element) ParentNode() Node {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.parent == nil {
		return nil
	}
	return e.parent
}

func (e *element) FirstChild() Node {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.firstChild == nil {
		return nil
	}
	return e.firstChild
}

func (e *element) LastChild() Node {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.lastChild == nil {
		return nil
	}
	return e.lastChild
}

func (e *element) PreviousSibling() Node {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.prevSibling == nil {
		return nil
	}
	return e.prevSibling
}

func (e *element) NextSibling() Node {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.nextSibling == nil {
		return nil
	}
	return e.nextSibling
}

func (e *element) InsertBefore(newNode, ref Node) Node {
	newCNode := UnWrap(newNode)
	refNode := UnWrap(ref)

	if refNode == nil {
		e.AppendChild(newNode)
		return nil
	}

	e.rw.Lock()
	defer e.rw.Unlock()

	if newChild, ok := newCNode.(*element); ok {
		if oldChild, ok := refNode.(*element); ok {
			if newChild.parent != nil || newChild.prevSibling != nil || newChild.nextSibling != nil {
				panic("dom: insertBefore called for an attached child Element")
			}
			var prev, next *element
			if oldChild != nil {
				prev, next = oldChild.prevSibling, oldChild
			} else {
				prev = e.lastChild
			}
			if prev != nil {
				prev.nextSibling = newChild
			} else {
				e.firstChild = newChild
			}
			if next != nil {
				next.prevSibling = newChild
			} else {
				e.lastChild = newChild
			}
			newChild.parent = e
			newChild.prevSibling = prev
			newChild.nextSibling = next

			return oldChild
		}
	}

	return nil
}

func (e *element) AppendChild(n Node) {
	e.rw.Lock()
	defer e.rw.Unlock()

	p := UnWrap(n)
	if p == nil {
		return
	}
	c, ok := p.(*element)
	if !ok {
		panic("dom: AppendChild append nil child")
	}

	if c.parent != nil || c.prevSibling != nil || c.nextSibling != nil {
		panic("dom: AppendChild called for an attached child Element")
	}
	last := e.lastChild
	if last != nil {
		last.nextSibling = c
	} else {
		e.firstChild = c
	}
	e.lastChild = c
	c.parent = e
	c.prevSibling = last
}

func (e *element) RemoveChild(n Node) Node {
	e.rw.Lock()
	defer e.rw.Unlock()

	p := UnWrap(n)

	if p == nil {
		return nil
	}
	c := p.(*element)

	if c.parent != e {
		panic("dom: removeChild called for a non-child Element")
	}
	if e.firstChild == c {
		e.firstChild = c.nextSibling
	}
	if c.nextSibling != nil {
		c.nextSibling.prevSibling = c.prevSibling
	}
	if e.lastChild == c {
		e.lastChild = c.prevSibling
	}
	if c.prevSibling != nil {
		c.prevSibling.nextSibling = c.nextSibling
	}
	c.parent = nil
	c.prevSibling = nil
	c.nextSibling = nil

	return c
}

func (e *element) SetAttribute(k string, v interface{}) {
	e.rw.Lock()
	defer e.rw.Unlock()

	if e.attributes == nil {
		e.attributes = map[string]interface{}{}
	}
	e.attributes[k] = v
}

func (e *element) GetAttribute(k string) interface{} {
	e.rw.RLock()
	defer e.rw.RUnlock()

	if e.attributes == nil {
		return ""
	}
	return e.attributes[k]
}

func (e *element) RemoveAttribute(k string) {
	e.rw.Lock()
	defer e.rw.Unlock()

	if e.attributes != nil {
		delete(e.attributes, k)
	}
}

func (e *element) GetAttributeNames() []string {
	e.rw.RLock()
	defer e.rw.RUnlock()

	names := make([]string, 0, len(e.attributes))

	for k := range e.attributes {
		names = append(names, k)
	}

	return names
}

func (*element) AddEventListener(eventName string, listener func(Event), args ...interface{}) {
	// todo
}

func (*element) RemoveEventListener(eventName string, listener func(Event), args ...interface{}) {
	// todo
}

func (*element) DispatchEvent(event Event) {
	// todo
}
