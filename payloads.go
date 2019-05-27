package mvds

import "sync"

type Payloads struct {
	sync.RWMutex

	payloads map[GroupID]map[PeerId]Payload
}

func (p *Payloads) Map(f func(GroupID, PeerId, Payload)) {
	p.RLock()
	defer p.RUnlock()

	for g, payloads := range p.payloads {
		for peer, payload := range payloads {
			f(g, peer, payload)
		}
	}
}

func (p *Payloads) RemoveAll() {
	p.Lock()
	defer p.Unlock()

	p.payloads = make(map[GroupID]map[PeerId]Payload)
}
