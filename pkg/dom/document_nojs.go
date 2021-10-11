//go:build !js
// +build !js

package dom

var (
	Document = &document{
		root: &element{
			nodeType: DOCUMENT_NODE,
		},
	}
)

type document struct {
	root Element

	object
}

func (document) QuerySelector(selector string) Element {
	// todo
	return nil
}

func (document) CreateElement(tagName string) Element {
	return &element{
		nodeType: ELEMENT_NODE,
		tagName:  tagName,
	}
}

func (document) CreateTextNode(data string) Element {
	return &element{
		nodeType:    TEXT_NODE,
		textContent: data,
	}
}
