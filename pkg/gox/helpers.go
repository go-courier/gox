package gox

import "reflect"

func RangeMap(v interface{}, m func(key string) interface{}) []interface{} {
	t := reflect.TypeOf(v)

	switch t.Kind() {
	case reflect.Map:
		rv := reflect.ValueOf(v)

		children := make([]interface{}, 0, rv.Len())

		iter := reflect.ValueOf(v).MapRange()
		for iter.Next() {
			children = append(children, m(iter.Key().String()))
		}

		return children
	}

	return nil
}

func RangeSlice(v interface{}, m func(i int) interface{}) []interface{} {
	t := reflect.TypeOf(v)

	switch t.Kind() {
	case reflect.Slice:
		rv := reflect.ValueOf(v)

		children := make([]interface{}, rv.Len())

		for i := range children {
			children[i] = m(i)
		}

		return children
	}

	return nil
}
