package mvds

import (
	"reflect"
	"sync"
)

type Payloads struct {
	sync.RWMutex

	payloads map[GroupID]map[PeerId]Payload
}

func (p *Payloads) get(id GroupID, peerId PeerId) Payload {
	payload, _ := p.payloads[id][peerId]
	return payload
}

func (p *Payloads) set(id GroupID, peerId PeerId, payload Payload) {
	_, ok := p.payloads[id]
	if !ok {
		p.payloads[id] = make(map[PeerId]Payload)
	}

	p.payloads[id][peerId] = payload
}

// @todo check in all the functions below that we aren't duplicating stuff

func (p *Payloads) AddOffers(group GroupID, peer PeerId, offers ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Offer == nil {
		payload.Offer = &Offer{Id: make([][]byte, 0)}
	}

	for _, o := range offers {
		if isInArray(payload.Offer.Id, o) {
			continue
		}

		payload.Offer.Id = append(payload.Offer.Id, o)
	}

	p.set(group, peer, payload)
}

func (p *Payloads) AddAcks(group GroupID, peer PeerId, acks ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Ack == nil {
		payload.Ack = &Ack{Id: make([][]byte, 0)}
	}

	for _, a := range acks {
		if isInArray(payload.Ack.Id, a) {
			continue
		}

		payload.Ack.Id = append(payload.Ack.Id, a)
	}

	p.set(group, peer, payload)
}

func (p *Payloads) AddRequests(group GroupID, peer PeerId, request ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Request == nil {
		payload.Request = &Request{Id: make([][]byte, 0)}
	}

	for _, r := range request {
		if isInArray(payload.Request.Id, r) {
			continue
		}

		payload.Request.Id = append(payload.Request.Id, r)
	}

	p.set(group, peer, payload)
}

func (p *Payloads) AddMessages(group GroupID, peer PeerId, messages ...*Message) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Messages == nil {
		payload.Messages = make([]*Message, 0)
	}

	payload.Messages = append(payload.Messages, messages...)
	p.set(group, peer, payload)
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

func isInArray(a [][]byte, val []byte) bool {
	for _, v := range a {
		if reflect.DeepEqual(v, val) {
			return true
		}
	}

	return false
}
