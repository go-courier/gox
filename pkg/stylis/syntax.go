package stylis

import (
	"bytes"
	"github.com/go-courier/jsx/pkg/stylis/scanner"
	"io"
	"strings"
)

type FormatOpt struct {
	Indent  string
	Depth   int
	OneLine bool
}

func (opt *FormatOpt) WriteIdent(w io.Writer) {
	for i := 0; i < opt.Depth; i++ {
		_, _ = io.WriteString(w, opt.Indent)
	}
}

func (opt *FormatOpt) WriteLine(w io.Writer, s string) {
	if !opt.OneLine {
		opt.WriteIdent(w)
	}

	_, _ = io.WriteString(w, s)

	if !opt.OneLine {
		_, _ = io.WriteString(w, "\n")
	}
}

type Node interface {
	IsEmpty() bool
	FormatTo(w io.Writer, opt *FormatOpt)
}

type NodeAppendable interface {
	AppendNode(n Node)
}

type WithSelector interface {
	WithSelector(parent string) []Node
}

type Root struct {
	Nodes []Node
}

func (r *Root) IsEmpty() bool {
	return len(r.Nodes) == 0
}

func (r *Root) AppendNode(n Node) {
	r.Nodes = append(r.Nodes, n)
}

func (r *Root) FormatTo(w io.Writer, opt *FormatOpt) {
	if opt == nil {
		opt = &FormatOpt{
			OneLine: false,
			Indent:  "  ",
		}
	}

	WalkRules(r.Nodes, "", func(n Node) {
		if !n.IsEmpty() {
			n.FormatTo(w, opt)
			_, _ = io.WriteString(w, "\n")
		}
	})
}

type AtRule struct {
	Name   string
	Params string
	Nodes  []Node
}

func (r *AtRule) IsEmpty() bool {
	return len(r.Nodes) == 0
}

func (r *AtRule) AppendNode(n Node) {
	r.Nodes = append(r.Nodes, n)
}

func (r *AtRule) WithSelector(parent string) (rules []Node) {
	root := &AtRule{
		Name:   r.Name,
		Params: r.Params,
	}

	rules = append(rules, root)

	decls := make([]Node, 0, len(r.Nodes))

	for i := range r.Nodes {
		switch x := r.Nodes[i].(type) {
		case WithSelector:
			root.Nodes = append(root.Nodes, x.WithSelector(parent)...)
		case *Declaration:
			if root.Name == "media" {
				decls = append(decls, x)
			} else {
				root.Nodes = append(root.Nodes, x)
			}
		default:
			root.Nodes = append(root.Nodes, x)
		}
	}

	if len(decls) > 0 {
		root.Nodes = append(root.Nodes, (&Rule{
			Selector: "&",
			Nodes:    decls,
		}).WithSelector(parent)...)
	}

	return
}

func (r *AtRule) AtRule() string {
	if r.Params == "" {
		return r.Name
	}
	return r.Name + " " + r.Params
}

func (r *AtRule) FormatTo(w io.Writer, opt *FormatOpt) {
	if opt == nil {
		opt = &FormatOpt{
			OneLine: false,
			Indent:  "  ",
		}
	}

	opt.WriteLine(w, "@"+r.AtRule()+"{")

	o := *opt
	o.Depth++

	WalkRules(r.Nodes, "@", func(r Node) {
		if !r.IsEmpty() {
			r.FormatTo(w, &o)
		}
	})

	opt.WriteLine(w, "}")
}

type Rule struct {
	Selector string
	Nodes    []Node
}

func (r *Rule) IsEmpty() bool {
	return len(r.Nodes) == 0
}

func (r *Rule) AppendNode(n Node) {
	r.Nodes = append(r.Nodes, n)
}

func (r *Rule) WithSelector(parent string) (rules []Node) {
	root := &Rule{}

	if parent == "" {
		root.Selector = r.Selector
	} else {
		root.Selector = concatSelectors(
			strings.Split(r.Selector, ","),
			strings.Split(parent, ","),
		)
	}

	rules = append(rules, root)

	for i := range r.Nodes {
		switch x := r.Nodes[i].(type) {
		case WithSelector:
			rules = append(rules, x.WithSelector(root.Selector)...)
		default:
			root.Nodes = append(root.Nodes, x)
		}
	}

	return
}

func (r *Rule) FormatTo(w io.Writer, opt *FormatOpt) {
	if len(r.Nodes) == 0 {
		return
	}

	if opt == nil {
		opt = &FormatOpt{
			OneLine: false,
			Indent:  "  ",
		}
	}

	opt.WriteLine(w, r.Selector+"{")

	o := *opt
	o.Depth++

	WalkRules(r.Nodes, r.Selector, func(r Node) {
		r.FormatTo(w, &o)
	})

	opt.WriteLine(w, "}")
}

type Declaration struct {
	Prop  string
	Value string
}

func (s *Declaration) IsVariable() bool {
	return strings.HasPrefix(s.Prop, "--")
}

func (s *Declaration) IsEmpty() bool {
	return s.Value == ""
}

func (s *Declaration) FormatTo(w io.Writer, opt *FormatOpt) {
	opt.WriteLine(w, s.Prop+":"+s.Value+";")
}

func WalkRules(rules []Node, parent string, each func(r Node)) {
	for i := range rules {
		switch x := rules[i].(type) {
		case WithSelector:
			if parent != "@" {
				composedRules := x.WithSelector(parent)
				for i := range composedRules {
					each(composedRules[i])
				}
			} else {
				each(rules[i])
			}
		default:
			each(x)
		}
	}
}

func concatSelectors(selectors []string, parentSelectors []string) string {
	finals := make([]string, 0, len(selectors)*len(parentSelectors))

	for i := range selectors {
		s := strings.TrimSpace(selectors[i])

		for j := range parentSelectors {
			parent := strings.TrimSpace(parentSelectors[j])
			if parent != "" {
				finals = append(finals, ReplaceHolder(s, parent))
			}
		}
	}

	return strings.Join(finals, ",")
}

func ReplaceHolder(rule string, p string) string {
	data := []rune(rule)
	pos := 0

	next := func() rune {
		pos++
		if pos < len(data) {
			return data[pos]
		}
		return scanner.TokenEOF
	}

	start := 0

	buf := bytes.NewBuffer(nil)
	hasHolder := false

	for pos < len(data) {
		start = pos
		ch := data[pos]

		if ch == '&' {
			ch = next()
			hasHolder = true
			_, _ = buf.WriteString(p)
			continue
		}

		switch ch {
		case '\\':
			ch, _ = scanner.ScanEscape(ch, next)
			ch = next()
		case '"', '\'':
			ch, _ = scanner.ScanString(ch, next)
			ch = next()
		default:
			ch = next()
		}

		_, _ = buf.WriteString(string(data[start:pos]))
	}

	if hasHolder {
		return buf.String()
	}

	return p + " " + buf.String()
}
