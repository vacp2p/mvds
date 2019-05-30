package epoch

import "sync"

type Epoch struct {
	sync.RWMutex

	e uint64
}

func NewEpoch() *Epoch {
	return &Epoch{e: 0}
}

func (e *Epoch) Current() uint64 {
	e.RLock()
	defer e.RUnlock()

	return e.e
}

func (e *Epoch) Increment() {
	e.Lock()
	defer e.Unlock()

	e.e += 1
}
