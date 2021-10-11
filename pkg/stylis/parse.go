package stylis

import (
	"bytes"
	"github.com/go-courier/jsx/pkg/stylis/scanner"
	"strings"
)

func Parse(rule string) Node {
	return ParseBytes([]byte(rule))
}

func ParseBytes(rule []byte) Node {
	s := &nodeParser{s: scanner.NewCSSScanner(bytes.NewBuffer(rule))}
	return s.ScanNode()
}

type nodeParser struct {
	s           *scanner.Scanner
	blockStacks []NodeAppendable
}

func (p *nodeParser) ScanNode() Node {
	root := &Root{}

	p.blockStacks = []NodeAppendable{root}

	lastTokenText := ""

	for tok := p.s.Scan(); tok != scanner.TokenEOF; tok = p.s.Scan() {
		switch tok {
		case scanner.TokenSemicolon:
			p.addDecl(lastTokenText)
		case scanner.TokenBlockStart:
			p.openBlock(lastTokenText)
		case scanner.TokenBlockEnd:
			if lastTokenText != ";" {
				p.addDecl(lastTokenText)
			}
			p.closeBlock()
		}
		lastTokenText = p.s.TokenText()
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
