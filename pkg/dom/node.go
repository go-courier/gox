package dom

type NodeType uint32

// https://developer.mozilla.org/en-US/docs/Web/API/Node
const (
	ERROR_NODE NodeType = iota
	ELEMENT_NODE
	ATTRIBUTE_NODE
	TEXT_NODE
	CDATA_SECTION_NODE
	ENTITY_REFERENCE_NODE
	ENTITY_NODE
	PROCESSING_INSTRUCTION_NODE
	COMMENT_NODE
	DOCUMENT_NODE
	DOCUMENT_TYPE_NODE
	DOCUMENT_FRAGMENT_NODE
	NOTATION_NODE
)

// Node https://developer.mozilla.org/en-US/docs/Web/API/Node
type Node interface {
	JSObject

	EventTarget

	NodeType() NodeType
	NodeName() string

	OwnerDocument() Doc

	ParentNode() Node
	FirstChild() Node
	LastChild() Node
	PreviousSibling() Node
	NextSibling() Node

	TextContent() string
	SetTextContent(textContent string)

	InsertBefore(newChild, referenceNode Node) Node
	RemoveChild(n Node) Node
	AppendChild(n Node)
}

func UnWrap(v interface{}) interface{} {
	if can, ok := v.(interface{ Underlying() interface{} }); ok {
		return UnWrap(can.Underlying())
	}
	return v
}
