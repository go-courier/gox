package stylis

import (
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
	return IsNodesEmpty(r.Nodes)
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
	return IsNodesEmpty(r.Nodes)
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
		n := r.Nodes[i]
		switch x := n.(type) {
		case WithSelector:
			switch root.Name {
			case "-webkit-keyframes", "keyframes":
				root.Nodes = append(root.Nodes, n)
			default:
				root.Nodes = append(root.Nodes, x.WithSelector(parent)...)
			}
		case *Declaration:
			switch root.Name {
			case "supports":
				if x.IsVariable() {
					root.Nodes = append(root.Nodes, x)
				} else {
					decls = append(decls, x)
				}
			case "media":
				decls = append(decls, x)
			default:
				root.Nodes = append(root.Nodes, x)

			}
		default:
			root.Nodes = append(root.Nodes, x)
		}
	}

	if len(decls) > 0 {
		root.Nodes = append((&Rule{
			Selector: "&",
			Nodes:    decls,
		}).WithSelector(parent), root.Nodes...)
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
		r.FormatTo(w, &o)
	})

	opt.WriteLine(w, "}")
}

type Rule struct {
	Selector string
	Nodes    []Node
}

func (r *Rule) IsEmpty() bool {
	return IsNodesEmpty(r.Nodes)
}

func IsNodesEmpty(nodes []Node) bool {
	if len(nodes) == 0 {
		return true
	}

	for i := range nodes {
		if !nodes[i].IsEmpty() {
			return false
		}
	}
	return true
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
			SplitRules(r.Selector),
			SplitRules(parent),
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
					r := composedRules[i]
					if r.IsEmpty() {
						continue
					}
					each(composedRules[i])
				}
			} else {
				r := rules[i]
				if r.IsEmpty() {
					continue
				}
				each(r)
			}
		default:
			if x.IsEmpty() {
				continue
			}
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
