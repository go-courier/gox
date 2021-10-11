package htmlloader

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/go-courier/gox/pkg/pwasm/core"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type HTMLLoader struct{}

func isHttpHref(href string) bool {
	return strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") || strings.HasPrefix(href, "//")
}

func (HTMLLoader) Transform(ctx context.Context, filename string, src io.Reader) (o core.Output, err error) {
	doc, err := goquery.NewDocumentFromReader(src)
	if err != nil {
		return nil, err
	}

	p := core.ProcessorFromContext(ctx)

	o = &OutputHTML{
		filename: filename,
		doc:      doc,
	}

	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			if isHttpHref(href) {
				return
			}
			resolvedPath, e := p.Resolve(ctx, href, filename, false)
			if e != nil {
				err = e
				return
			}
			s.SetAttr("href", resolvedPath)
		}
	})

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if code := s.Text(); code != "" {
			output, e := p.RoundTrip(ctx, core.NewSourceFile(filename+"~index.ts", []byte(code), false))
			if e != nil {
				err = e
				return
			}
			s.SetHtml(string(output.Code()))
		} else {
			if href, ok := s.Attr("src"); ok {
				if isHttpHref(href) {
					return
				}

				switch role, _ := s.Attr("data-role"); role {
				case "service-worker":
					s.RemoveAttr("src")

					resolvedPath, e := p.Resolve(ctx, href, filename, true)
					if e != nil {
						err = e
						return
					}

					moduleType := func() string {
						if t, ok := s.Attr("type"); ok {
							return strconv.Quote(t)
						}
						return "undefined"
					}

					output, e := p.RoundTrip(ctx, core.NewSourceFile(filename+"~"+href, []byte(`
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('`+resolvedPath+`', { type: `+moduleType()+`});
  });
}
`), false))
					if e != nil {
						err = e
						return
					}
					s.SetHtml(string(output.Code()))
				default:
					resolvedPath, e := p.Resolve(ctx, href, filename, false)
					if e != nil {
						err = e
						return
					}
					s.SetAttr("src", resolvedPath)
				}
			}
		}
	})

	return
}

type OutputHTML struct {
	filename string
	doc      *goquery.Document
}

func (o *OutputHTML) Filename() string {
	return o.filename
}

func (o *OutputHTML) Code() []byte {
	buf := bytes.NewBuffer(nil)
	for i := range o.doc.Nodes {
		n := o.doc.Nodes[i]
		_ = html.Render(buf, n)
	}
	return buf.Bytes()
}
