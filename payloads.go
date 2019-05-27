package mvds

import "sync"

type Payloads struct {
	sync.RWMutex

	payloads map[GroupID]map[PeerId]Payload
}

func (p *Payloads) Map(f func(GroupID, PeerId, Payload)) {
	for g, payloads := range p.payloads {
		for peer, payload := range payloads {
			f(g, peer, payload)
		}
	}
}
