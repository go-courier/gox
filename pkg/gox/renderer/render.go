package renderer

import (
	"context"
	"io"

	. "github.com/go-courier/gox/pkg/dom"
	"github.com/go-courier/gox/pkg/gox"
)

func RenderToString(ctx context.Context, w io.Writer, vnode *gox.VNode) error {
	doc := Document.CreateElement("html")

	if err := CreateRoot(doc).Render(ctx, vnode); err != nil {
		return err
	}

	RenderToHTML(w, doc)
	return nil
}
