package dom

type JSObject interface {
	Get(propName string) interface{}
	Set(propName string, v interface{})
}

type object struct {
	props map[string]interface{}
}

func (e *object) Get(name string) interface{} {
	if e.props == nil {
		return nil
	}
	if v, ok := e.props[name]; ok {
		return v
	}
	return nil
}

func (e *object) Set(name string, value interface{}) {
	if e.props == nil {
		e.props = map[string]interface{}{}
	}
	e.props[name] = value
}
