package internal

import "context"

type Component interface {
	Render(ctx context.Context, children ...interface{}) interface{}
}
