package internal

import "context"

type CanAttrs interface {
	Attrs(ctx context.Context) Attrs
}

type Attrs map[string]interface{}

func (attrs Attrs) Merge(attrsList ...Attrs) {
	for k := range attrsList {
		for k, vv := range attrsList[k] {
			attrs[k] = vv
		}
	}
}

func (attrs Attrs) Each(fn func(k string, v interface{})) {
	if len(attrs) > 0 {
		for k := range attrs {
			fn(k, attrs[k])
		}
	}
}
