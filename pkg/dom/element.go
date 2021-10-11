package dom

// Element https://developer.mozilla.org/en-US/docs/Web/API/Element
type Element interface {
	Node

	GetAttributeNames() []string
	SetAttribute(k string, v interface{})
	GetAttribute(k string) interface{}
	RemoveAttribute(k string)
}

type ElementList []Element

func (list ElementList) Each(fn func(n Element)) {
	for i := range list {
		fn(list[i])
	}
}
