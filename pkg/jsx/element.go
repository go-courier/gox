package jsx

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/go-courier/jsx/pkg/css"
)

type CSS = css.CSS

type Element interface {
	RenderTo(ctx context.Context, w io.Writer) error
}

func T(tag string) func(children ...Element) Component {
	return func(children ...Element) Component {
		return func(ctx context.Context) Element {
			styles := make([]CSS, 0, len(children))
			attributes := make([]Attributes, 0, len(children))
			finals := make([]Element, 0, len(children))

			for i := range children {
				c := children[i]

				switch x := c.(type) {
				case Attributes:
					attributes = append(attributes, x)
				case CSS:
					styles = append(styles, x)
				default:
					finals = append(finals, x)
				}
			}

			if len(styles) > 0 {
				c := css.CacheFromContext(ctx)
				if c != nil {
					attributes = append(attributes, Attr("class", c.CSS(styles...)))
				} else {
					attributes = append(attributes, Attr("style", css.MergeCSS(styles...)))
				}
			}

			return &elem{
				tag:      tag,
				attrs:    MergeAttrs(attributes...),
				children: finals,
			}
		}

	}
}

type elem struct {
	tag      string
	attrs    map[string]interface{}
	children []Element
}

func (e *elem) RenderTo(ctx context.Context, w io.Writer) error {
	_, _ = fmt.Fprintf(w, "<%s", e.tag)

	if len(e.attrs) > 0 {
		attrKeys := make([]string, 0)

		for k := range e.attrs {
			attrKeys = append(attrKeys, k)
		}

		sort.Strings(attrKeys)

		for _, k := range attrKeys {
			switch v := e.attrs[k].(type) {
			case bool:
				if v {
					_, _ = fmt.Fprintf(w, " %s", k)
				}
			case fmt.Stringer:
				_, _ = fmt.Fprintf(w, " %s=%s", k, strconv.Quote(v.String()))
			default:
				_, _ = fmt.Fprintf(w, " %s=%s", k, strconv.Quote(fmt.Sprintf("%v", v)))
			}
		}
	}

	_, _ = fmt.Fprintf(w, ">")

	if len(e.children) > 0 {
		if err := render(ctx, Fragment(e.children), w); err != nil {
			return err
		}
	}

	_, _ = fmt.Fprintf(w, "</%s>", e.tag)

	return nil
}
