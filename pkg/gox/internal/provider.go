package internal

import "context"

type ContextProvider interface {
	GetChildContext(ctx context.Context) context.Context
}

type Provider func(ctx context.Context) context.Context

func (Provider) Compare(j interface{}) int {
	if _, ok := j.(Provider); ok {
		return 0
	}
	return 1
}

func (p Provider) GetChildContext(ctx context.Context) context.Context {
	return p(ctx)
}

func (Provider) Render(ctx context.Context, children ...interface{}) interface{} {
	return JSX(Fragment{}, children...)
}
