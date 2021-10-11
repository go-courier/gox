package internal

import "context"

type Text string

func (Text) builtin() {
}

func (Text) Render(ctx context.Context, children ...interface{}) interface{} {
	return nil
}
