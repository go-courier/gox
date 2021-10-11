package internal

type Hook interface {
	Update(next Hook)
}

type HookCanCommit interface {
	Commit()
}

type HookCanDestroy interface {
	Destroy()
}

type hooks struct {
	locked     bool
	usedHooks  []Hook
	hookUseIdx int
}

func (hs *hooks) hookInit() {
	hs.hookUseIdx = 0
}

func (hs *hooks) putHook(hook Hook, idx int) Hook {
	if maxIdx := len(hs.usedHooks) - 1; maxIdx < hs.hookUseIdx {
		hs.usedHooks = append(hs.usedHooks, hook)
		return hs.usedHooks[idx]
	}
	return hs.usedHooks[idx]
}

func (hs *hooks) use(hook Hook) Hook {
	usedHook := hs.putHook(hook, hs.hookUseIdx)
	usedHook.Update(hook)
	hs.hookUseIdx++
	return usedHook
}

func (hs *hooks) commit() {
	for i := range hs.usedHooks {
		if hc, ok := hs.usedHooks[i].(HookCanCommit); ok {
			hc.Commit()
		}
	}
}

func (hs *hooks) destroy() {
	for i := range hs.usedHooks {
		if hc, ok := hs.usedHooks[i].(HookCanDestroy); ok {
			hc.Destroy()
		}
	}
}
