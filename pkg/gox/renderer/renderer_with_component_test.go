package renderer_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/go-courier/gox/pkg/gox/renderer"

	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
	"github.com/onsi/gomega"
)

type AppSubComp struct {
	Text string
}

func (a AppSubComp) Render(ctx context.Context, children ...interface{}) interface{} {
	return a.Text
}

type App struct {
	Text string
}

func (a App) Render(ctx context.Context, children ...interface{}) interface{} {
	return Div("hello ", H(AppSubComp{a.Text})())
}

func TestRenderWithComponent(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := context.Background()
	root := Document.CreateElement("body")
	r := renderer.CreateRoot(root)

	_ = r.Render(ctx, H(App{Text: "app"})())

	buf.Reset()
	RenderToHTML(buf, root)
	gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div>hello app</div></body>`))

	t.Run("should update when props changed", func(t *testing.T) {
		_ = r.Render(ctx, H(App{Text: "app updated"})())

		buf.Reset()
		RenderToHTML(buf, root)
		gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div>hello app updated</div></body>`))
	})
}
