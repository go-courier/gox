package renderer_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
	"github.com/go-courier/gox/pkg/gox/renderer"
	"github.com/onsi/gomega"
)

func TestRender(t *testing.T) {
	t.Run("should update vnode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := context.Background()
		root := Document.CreateElement("body")
		r := renderer.CreateRoot(root)

		{
			buf.Reset()
			_ = r.Render(ctx, Div(Attrs{"role": "value"}, Span("1")))
			fmt.Println("render")
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div role="value"><span>1</span></div></body>`))
		}

		{
			buf.Reset()
			_ = r.Render(ctx, Div(Attrs{"role": "value"}, Span("2")))
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div role="value"><span>2</span></div></body>`))
		}

		{
			buf.Reset()
			_ = r.Render(ctx, Div(Attrs{"role": "value", "id": "test"}, Fragment("2")))
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div id="test" role="value">2</div></body>`))
		}

		{
			buf.Reset()
			_ = r.Render(ctx, Div(Attrs{"role": "value"}, "1"))
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div role="value">1</div></body>`))
		}

		{
			buf.Reset()
			_ = r.Render(ctx, Span("1"))
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><span>1</span></body>`))
		}
	})

	t.Run("should render to portal", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		ctx := context.Background()
		root := Document.CreateElement("body")
		portal := Document.CreateElement("portal")
		r := renderer.CreateRoot(root)

		{
			_ = r.Render(ctx, Div(
				Span("1"),
				renderer.Portal(portal)(Span("2")),
			))

			buf.Reset()
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div><span>1</span></div></body>`))

			buf.Reset()
			RenderToHTML(buf, portal)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<portal><span>2</span></portal>`))
		}

		{
			_ = r.Render(ctx, Div(
				Span("11"),
				renderer.Portal(portal)(Span("22")),
			))

			buf.Reset()
			RenderToHTML(buf, root)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div><span>11</span></div></body>`))

			buf.Reset()
			RenderToHTML(buf, portal)
			gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<portal><span>22</span></portal>`))
		}
	})
}
