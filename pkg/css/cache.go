package css

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/go-courier/gox/pkg/dom"
	"github.com/go-courier/gox/pkg/hash"
	"github.com/go-courier/gox/pkg/stylis"
)

type SerializedStyles struct {
	Name   string
	Styles []byte
}

func (s *SerializedStyles) ToCSS(key string) []byte {
	b := bytes.NewBuffer(nil)

	n := stylis.Parse(fmt.Sprintf(".%s-%s{%s}", key, s.Name, string(s.Styles)))

	n.FormatTo(b, &stylis.FormatOpt{
		OneLine: true,
	})

	return bytes.TrimSpace(b.Bytes())
}

func NewCSSCache(key string, mount dom.Element) *CSSCache {
	return &CSSCache{
		Key:   key,
		mount: mount,
	}
}

type CSSCache struct {
	Key        string
	Registered map[string]*SerializedStyles
	mount      dom.Element
}

func (c *CSSCache) CSS(ctx context.Context, csses ...CSS) string {
	s := c.SerializeStyles(ctx, MergeCSS(csses...))
	return c.Key + "-" + s.Name
}

func (c *CSSCache) SerializeStyles(ctx context.Context, args ...interface{}) *SerializedStyles {
	ss := &SerializedStyles{}

	for i := range args {
		switch x := args[i].(type) {
		case CSS:
			ss.Styles = append(ss.Styles, toStylesBytes(x)...)
		}
	}

	ss.Name = hash.MurmurHash2String(ss.Styles, 0)

	if c.Registered == nil {
		c.Registered = map[string]*SerializedStyles{}
	}

	if ss, ok := c.Registered[ss.Name]; ok {
		return ss
	}

	c.Mount(ctx, ss)
	c.Registered[ss.Name] = ss

	return ss
}

func (c *CSSCache) Mount(ctx context.Context, ss *SerializedStyles) {
	if c.mount == nil {
		return
	}

	doc := c.mount.OwnerDocument()

	d := doc.CreateElement("style")
	d.SetAttribute("data-css", ss.Name)
	d.AppendChild(doc.CreateTextNode(string(ss.ToCSS(c.Key))))

	c.mount.AppendChild(d)
}

func toStylesBytes(s map[string]interface{}) []byte {
	if len(s) == 0 {
		return nil
	}

	keys := make([]string, 0)

	for k := range s {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	b := bytes.NewBuffer(nil)

	for _, k := range keys {
		v := s[k]
		_, _ = fmt.Fprintf(b, "%s", toSnake(k))

		switch x := v.(type) {
		case CSS:
			b.WriteByte('{')
			b.Write(toStylesBytes(x))
			b.WriteByte('}')
		case fmt.Stringer:
			_, _ = fmt.Fprintf(b, ":%s", x.String())
		default:
			_, _ = fmt.Fprintf(b, ":%v", v)
		}

		b.WriteByte(';')
	}

	return b.Bytes()
}

func toSnake(id string) string {
	buf := bytes.NewBuffer(nil)

	chars := []byte(id)

	for _, c := range chars {
		if c >= 'A' && c <= 'Z' {
			buf.WriteByte('-')
			c += 'a' - 'A'
		}
		buf.WriteByte(c)
	}

	return buf.String()
}

type contextKeyCSSCache struct{}

func CSSCacheFromContext(ctx context.Context) *CSSCache {
	if c, ok := ctx.Value(contextKeyCSSCache{}).(*CSSCache); ok {
		return c
	}
	return nil
}

func ContextWithCSSCache(ctx context.Context, c *CSSCache) context.Context {
	return context.WithValue(ctx, contextKeyCSSCache{}, c)
}
