package core

import (
	"context"
	"io"
)

type Loader interface {
	Transform(ctx context.Context, filename string, src io.Reader) (Output, error)
}

type Output interface {
	Filename() string
	Code() []byte
}
