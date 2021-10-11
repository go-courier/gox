package css

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/go-courier/jsx/pkg/hash"
	"github.com/go-courier/jsx/pkg/postcss"
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

func NewCache() *Cache {
	return &Cache{
		Key: "css",
	}
}

type Cache struct {
	Key        string
	Registered map[string]*SerializedStyles
}

func (c *Cache) CSS(csses ...CSS) string {
	s := c.SerializeStyles(MergeCSS(csses...))
	return c.Key + "-" + s.Name
}

func (c *Cache) SerializeStyles(args ...interface{}) *SerializedStyles {
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

	c.Registered[ss.Name] = ss

	return ss
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

type contextKeyCache struct{}

func CacheFromContext(ctx context.Context) *Cache {
	if c, ok := ctx.Value(contextKeyCache{}).(*Cache); ok {
		return c
	}
	return nil
}

func ContextWithCache(ctx context.Context, c *Cache) context.Context {
	return context.WithValue(ctx, contextKeyCache{}, c)
}
