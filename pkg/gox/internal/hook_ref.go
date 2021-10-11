package internal

import (
	"fmt"
)

type RefHook struct {
	Ref
}

func (RefHook) Update(next Hook) {
}

func (s *RefHook) String() string {
	return fmt.Sprintf("UseRef: %v", s.Ref.Current)
}
