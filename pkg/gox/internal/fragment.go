package internal

import "context"

type Fragment struct{}

func (Fragment) builtin() {}

func (Fragment) Render(ctx context.Context, children ...interface{}) interface{} {
	return nil
}
