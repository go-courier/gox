package internal

import (
	"fmt"
)

type MemoHook struct {
	Setup     func() interface{}
	Deps      []interface{}
	memorised interface{}
}

func (h *MemoHook) Memorised() interface{} {
	return h.memorised
}

func (h *MemoHook) String() string {
	return fmt.Sprintf("UseMemo: %v", h.Deps)
}

func (h *MemoHook) Destroy() {
	h.memorised = nil
}

func (h *MemoHook) Update(next Hook) {
	if nextHook, ok := next.(*MemoHook); ok {
		if nextHook == h || nextHook.Deps == nil || !ShallowEqual(nextHook.Deps, h.Deps) {
			h.memorised = nextHook.Setup()
		}
	}

}
