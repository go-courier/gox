package jsx

import (
	"context"
	"io"
)

type Component func(ctx context.Context) Element

func (fn Component) RenderTo(ctx context.Context, w io.Writer) error {
	return render(ctx, fn(ctx), w)
}
