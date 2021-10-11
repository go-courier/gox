package renderer_test

import (
	"bytes"
	"context"
	"testing"

	. "github.com/go-courier/gox/pkg/css"
	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
	"github.com/go-courier/gox/pkg/gox/renderer"
	"github.com/onsi/gomega"
)

type AppWithContext struct{}

func (a AppWithContext) Render(ctx context.Context, children ...interface{}) interface{} {
	css := NewCSSCache("app", nil)

	withCSSCache := func(c *CSSCache) func(ctx context.Context) context.Context {
		return func(ctx context.Context) context.Context {
			return ContextWithCSSCache(ctx, c)
		}
	}

	return Provider(withCSSCache(css))(
		Div(
			CSS{
				"color": "red",
			},
			1,
		),
	)
}

func TestRenderWithContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := context.Background()
	root := Document.CreateElement("body")
	r := renderer.CreateRoot(root)

	{
		buf.Reset()
		_ = r.Render(ctx, JSX(AppWithContext{}))
		RenderToHTML(buf, root)
		gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div class="app-tokvmb">1</div></body>`))
	}
}
