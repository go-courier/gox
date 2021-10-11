package dom

// EventTarget https://developer.mozilla.org/en-US/docs/Web/API/EventTarget
type EventTarget interface {
	AddEventListener(eventName string, listener func(event Event), args ...interface{})
	RemoveEventListener(eventName string, listener func(event Event), args ...interface{})
	DispatchEvent(event Event)
}

// Event https://developer.mozilla.org/en-US/docs/Web/API/Event
type Event interface {
	JSObject

	Target() EventTarget
	PreventDefault()
	StopImmediatePropagation()
	StopPropagation()
}
