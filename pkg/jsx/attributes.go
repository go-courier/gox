package jsx

import (
	"context"
	"io"
)

func MergeAttrs(props ...Attributes) Attributes {
	p := Attributes{}
	for _, item := range props {
		for k, vv := range item {
			p[k] = vv
		}
	}
	return p
}

type Attributes map[string]interface{}

func (Attributes) RenderTo(ctx context.Context, w io.Writer) error {
	return nil
}

func Attr(k string, v interface{}) Attributes {
	return Attributes{k: v}
}
