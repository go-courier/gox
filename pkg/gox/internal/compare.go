package internal

import "reflect"

type Comparer interface {
	Compare(v interface{}) int
}

func ShallowEqual(a interface{}, b interface{}) bool {
	switch xa := a.(type) {
	case Comparer:
		return xa.Compare(b) == 0
	case map[string]interface{}:
		switch xb := b.(type) {
		case map[string]interface{}:
			if len(xa) != len(xb) {
				return false
			}

			for k := range xa {
				if xbv, ok := xb[k]; ok {
					if xa[k] != xbv {
						return false
					}
				} else {
					return false
				}
			}

			return true
		}
		return false
	case []interface{}:
		switch xb := b.(type) {
		case []interface{}:
			if len(xa) != len(xb) {
				return false
			}

			for i := range xa {
				if !ShallowEqual(xa[i], xb[i]) {
					return false
				}
			}
			return true
		}
		return false
	}

	if a == nil || b == nil {
		return false
	}

	if reflect.TypeOf(a).Comparable() && reflect.TypeOf(b).Comparable() {
		return a == b
	}

	return false
}
