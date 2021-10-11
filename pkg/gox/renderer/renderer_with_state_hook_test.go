package renderer_test

import (
	"bytes"
	"context"
	"testing"

	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
	"github.com/go-courier/gox/pkg/gox/renderer"
	"github.com/onsi/gomega"
)

type AppWithStateHook struct {
	Value       string
	UpdateValue **func(v interface{})
}

func (a AppWithStateHook) Render(ctx context.Context, children ...interface{}) interface{} {
	value, updateValue := UseState(ctx, "hook")

	*a.UpdateValue = &updateValue

	return Span(a.Value, " ", JSX(AppWithStateHookSub{Value: value.(string)}))
}

type AppWithStateHookSub struct {
	Value string
}

func (a AppWithStateHookSub) Render(ctx context.Context, children ...interface{}) interface{} {
	return Fragment(a.Value)
}

func TestRenderWithStateHook(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	ctx := context.Background()
	root := Document.CreateElement("body")
	r := renderer.CreateRoot(root)

	t.Run("should re render when stage changed", func(t *testing.T) {
		var updateValue *func(v interface{})

		_ = r.Render(ctx, Div(
			H(AppWithStateHook{
				Value:       "hello",
				UpdateValue: &updateValue,
			})(),
		))

		buf.Reset()
		RenderToHTML(buf, root)
		gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div><span>hello hook</span></div></body>`))

		_ = r.Render(ctx, Div(
			H(AppWithStateHook{
				Value:       "hello2",
				UpdateValue: &updateValue,
			})(),
		))

		r.Act(func() {
			(*updateValue)("updated")
		})

		buf.Reset()
		RenderToHTML(buf, root)
		gomega.NewWithT(t).Expect(buf.String()).To(gomega.Equal(`<body><div><span>hello2 updated</span></div></body>`))
	})
}
