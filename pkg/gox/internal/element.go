package internal

import "context"

type Element string

func (Element) builtin() {
}

func (Element) Render(ctx context.Context, children ...interface{}) interface{} {
	return nil
}
