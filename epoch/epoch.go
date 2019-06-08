package epoch

import (
	"sync/atomic"
)

type Epoch uint64

func (e Epoch) Current() uint64 {
	return uint64(e)
}

func (e *Epoch) Increment() {
	*e = Epoch(atomic.AddUint64((*uint64)(e), 1))
}
