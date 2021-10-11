package stylis

import (
	"bytes"
	"strings"
)

func Parse(rule string) Node {
	return ParseBytes([]byte(rule))
}

func ParseBytes(rule []byte) Node {
	return (&nodeParser{}).ScanNode(rule)
}

type nodeParser struct {
	t           Tokenizer
	blockStacks []NodeAppendable
}

func (p *nodeParser) ScanNode(rule []byte) Node {
	root := &Root{}

	p.blockStacks = []NodeAppendable{root}

	buf := bytes.NewBuffer(nil)

	for _, token := range p.t.Tokenize(rule) {
		switch token {
		case ";":
			if lastTokenText := buf.String(); lastTokenText != "" {
				p.addDecl(lastTokenText)
			}
			buf.Reset()
		case "{":
			p.openBlock(buf.String())
			buf.Reset()
		case "}":
			if lastTokenText := buf.String(); lastTokenText != "" {
				p.addDecl(lastTokenText)
			}
			p.closeBlock()
			buf.Reset()
		default:
			buf.WriteString(token)
		}
	}

	return root
}

func (p *nodeParser) appendNode(n Node) {
	p.blockStacks[len(p.blockStacks)-1].AppendNode(n)

	if nodeAppendable, ok := n.(NodeAppendable); ok {
		p.blockStacks = append(p.blockStacks, nodeAppendable)
	}
}

func (p *nodeParser) closeBlock() {
	p.blockStacks = p.blockStacks[0 : len(p.blockStacks)-1]
}

func (p *nodeParser) openBlock(selector string) {
	selector = strings.TrimSpace(selector)

	if selector[0] == '@' {
		parts := strings.SplitN(selector[1:], " ", 2)

		p.appendNode(&AtRule{
			Name: parts[0],
			Params: func() string {
				if len(parts) > 1 {
					return parts[1]
				}
				return ""
			}(),
		})
	} else {
		p.appendNode(&Rule{
			Selector: selector,
		})
	}
}

func (p *nodeParser) addDecl(decl string) {
	parts := strings.SplitN(decl, ":", 2)

	if len(parts) != 2 {
		return
	}

	p.appendNode(&Declaration{
		Prop:  strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(parts[1]),
	})
}
