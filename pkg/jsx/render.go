package jsx

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"io"
	"net/http"
)

func RenderToResponse(ctx context.Context, rw http.ResponseWriter, e Element) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	rw.WriteHeader(http.StatusOK)
	_ = Render(ctx, e, rw)
}

func Render(ctx context.Context, e Element, w io.Writer) error {
	dw := &DeferWriter{
		list: list.New(),
	}

	if err := render(ctx, e, dw); err != nil {
		return err
	}

	_, err := dw.WriteTo(w)
	return err
}

type DeferWriter struct {
	buf  *bytes.Buffer
	list *list.List
}

func (dw *DeferWriter) WriteTo(w io.Writer) (int64, error) {
	for e := dw.list.Front(); e != nil; e = e.Next() {
		switch x := e.Value.(type) {
		case *bytes.Buffer:
			if _, err := io.Copy(w, x); err != nil {
				return -1, err
			}
		case Defer:
			if err := render(context.Background(), x(), w); err != nil {
				return -1, err
			}
		}
	}

	if dw.buf != nil {
		if _, err := io.Copy(w, dw.buf); err != nil {
			return -1, err
		}
	}

	return -1, nil
}

func (dw *DeferWriter) Write(p []byte) (n int, err error) {
	if dw.buf == nil {
		dw.buf = bytes.NewBuffer(nil)
	}
	return dw.buf.Write(p)
}

func (dw *DeferWriter) Defer(d Defer) {
	buf := dw.buf
	dw.list.PushBack(buf)
	dw.list.PushBack(d)
	dw.buf = nil
}

type Defer func() Element

func (d Defer) RenderTo(ctx context.Context, w io.Writer) error {
	if canDefer, ok := w.(interface{ Defer(Defer) }); ok {
		canDefer.Defer(d)
	}
	return nil
}

func render(ctx context.Context, e Element, w io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("render %T failed: %+v", e, r)
		}
	}()
	if e != nil {
		return e.RenderTo(ctx, w)
	}
	return nil
}
