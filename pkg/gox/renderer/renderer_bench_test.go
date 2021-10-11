package renderer

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-courier/gox/pkg/gox/internal"

	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
)

func Benchmark(b *testing.B) {
	b.Run("Single Node", func(b *testing.B) {
		b.Run("JSX", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = JSX(internal.Element(""))
			}
		})

		b.Run("mount", func(b *testing.B) {
			ctx := context.Background()
			root := Document.CreateElement("body")
			r := CreateRoot(root)

			d := Div()

			for i := 0; i < b.N; i++ {
				r.mount(ctx, d, Div())
			}
		})

		b.Run("patchVNode", func(b *testing.B) {
			ctx := context.Background()
			root := Document.CreateElement("body")
			r := CreateRoot(root)

			d := Div()

			for i := 0; i < b.N; i++ {
				r.patchVNode(ctx, d, Div())
			}
		})
	})

	b.Run("Large Nodes", func(b *testing.B) {
		b.Run("patchVNode", func(b *testing.B) {
			ctx := context.Background()
			root := Document.CreateElement("body")
			r := CreateRoot(root)

			children := make([]interface{}, 1000)

			for i := range children {
				children[i] = Span(fmt.Sprintf("%d", i))
			}

			for i := 0; i < b.N; i++ {
				_ = r.Render(ctx, Div(children...))
			}
		})
	})
}
