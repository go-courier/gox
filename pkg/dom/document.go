package dom

type Doc interface {
	JSObject

	CreateElement(tagName string) Element
	CreateTextNode(tagName string) Element
	QuerySelector(id string) Element
}
