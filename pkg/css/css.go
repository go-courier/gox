package css

import (
	"context"
	"io"
)

type CSS map[string]interface{}

func (CSS) RenderTo(ctx context.Context, w io.Writer) error {
	return nil
}

func MergeCSS(props ...CSS) CSS {
	p := CSS{}
	for _, item := range props {
		for k, vv := range item {
			p[k] = vv
		}
	}
	return p
}

func (s CSS) String() string {
	return string(toStylesBytes(s))
}
