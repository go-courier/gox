package core

import "context"

type processorContextKey struct{}

func ContextWithProcessor(ctx context.Context, p Processor) context.Context {
	return context.WithValue(ctx, processorContextKey{}, p)
}

func ProcessorFromContext(ctx context.Context) Processor {
	return ctx.Value(processorContextKey{}).(Processor)
}

type Processor interface {
	Resolve(ctx context.Context, filepath string, at string, async bool) (string, error)
	RoundTrip(ctx context.Context, output Output) (Output, error)
	Manifests() []Manifest
}

type Manifest struct {
	URL        string      `json:"url"`
	SourceFile *SourceFile `json:"-"`
}
