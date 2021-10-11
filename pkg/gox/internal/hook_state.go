package internal

import (
	"fmt"
)

type StateHook struct {
	State         interface{}
	OnStateChange func()
}

func (s *StateHook) Update(next Hook) {
	if n, ok := next.(*StateHook); ok {
		// context may change, should bind the latest callback
		s.OnStateChange = n.OnStateChange
	}
}

func (s *StateHook) String() string {
	return fmt.Sprintf("UseState: %v", s.State)
}

func (s *StateHook) SetState(nextStateOrValueFunc interface{}) {
	switch x := nextStateOrValueFunc.(type) {
	case func(v interface{}) interface{}:
		if nextState := x(s.State); nextState != s.State {
			s.State = nextState
			s.OnStateChange()
		}
	default:
		if x != s.State {
			s.State = x
			s.OnStateChange()
		}
	}
}
