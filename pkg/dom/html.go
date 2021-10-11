package dom

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func RenderToHTML(w io.Writer, n Element) {
	switch n.NodeType() {
	case TEXT_NODE:
		_, _ = io.WriteString(w, n.TextContent())
	case ELEMENT_NODE:
		_, _ = fmt.Fprintf(w, "<%s", n.NodeName())

		if names := n.GetAttributeNames(); len(names) > 0 {
			sort.Strings(names)

			for _, k := range names {
				_, _ = fmt.Fprintf(w, " %s=%s", k, strconv.Quote(stringify(n.GetAttribute(k))))
			}
		}

		_, _ = io.WriteString(w, ">")

		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			RenderToHTML(w, c.(Element))
		}

		_, _ = fmt.Fprintf(w, "</%s>", n.NodeName())
	}
}

func stringify(v interface{}) string {
	switch v := v.(type) {
	case []interface{}:
		s := &strings.Builder{}
		for i := range v {
			if i > 0 {
				s.WriteByte(' ')
			}
			s.WriteString(stringify(v[i]))
		}
		return s.String()
	case bool:
		return ""
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
