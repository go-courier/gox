package css

import (
	"context"

	"github.com/go-courier/gox/pkg/gox"
)

type CSS map[string]interface{}

func MergeCSS(props ...CSS) CSS {
	p := CSS{}
	for _, item := range props {
		for k, vv := range item {
			p[k] = vv
		}
	}
	return p
}

func (s CSS) Attrs(ctx context.Context) gox.Attrs {
	c := CSSCacheFromContext(ctx)

	if c != nil {
		return gox.Attrs{
			"class": c.CSS(ctx, s),
		}
	}

	return gox.Attrs{
		"style": s.Styles(),
	}
}

func (s CSS) Styles() string {
	return string(toStylesBytes(s))
}
