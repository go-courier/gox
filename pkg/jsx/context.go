package jsx

import (
	"context"
	"io"
)

func Provide(withContext func(ctx context.Context) context.Context) func(children ...Element) Element {
	return func(children ...Element) Element {
		return &ctx{withContext: withContext, children: children}
	}
}

type ctx struct {
	withContext func(ctx context.Context) context.Context
	children    []Element
}

func (c ctx) RenderTo(ctx context.Context, w io.Writer) error {
	ctx = c.withContext(ctx)

	for i := range c.children {
		if err := render(ctx, c.children[i], w); err != nil {
			return err
		}
	}

	return nil
}
