package mvds

import (
	"sync"
)

type payloads struct {
	sync.Mutex

	payloads map[GroupID]map[PeerId]Payload
}
// @todo check in all the functions below that we aren't duplicating stuff

func newPayloads() payloads {
	return payloads{
		payloads: make(map[GroupID]map[PeerId]Payload),
	}
}

func (p *payloads) AddOffers(group GroupID, peer PeerId, offers ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Offer == nil {
		payload.Offer = &Offer{Id: make([][]byte, 0)}
	}

	payload.Offer.Id = append(payload.Offer.Id, offers...)

	p.set(group, peer, payload)
}

func (p *payloads) AddAcks(group GroupID, peer PeerId, acks ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Ack == nil {
		payload.Ack = &Ack{Id: make([][]byte, 0)}
	}

	payload.Ack.Id = append(payload.Ack.Id, acks...)

	p.set(group, peer, payload)
}

func (p *payloads) AddRequests(group GroupID, peer PeerId, request ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Request == nil {
		payload.Request = &Request{Id: make([][]byte, 0)}
	}

	payload.Request.Id = append(payload.Request.Id, request...)

	p.set(group, peer, payload)
}

func (p *payloads) AddMessages(group GroupID, peer PeerId, messages ...*Message) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Messages == nil {
		payload.Messages = make([]*Message, 0)
	}

	payload.Messages = append(payload.Messages, messages...)
	p.set(group, peer, payload)
}

func (p *payloads) MapAndClear(f func(GroupID, PeerId, Payload)) {
	p.Lock()
	defer p.Unlock()

	for g, payloads := range p.payloads {
		for peer, payload := range payloads {
			f(g, peer, payload)
		}
	}

	p.payloads = make(map[GroupID]map[PeerId]Payload)
}

func (p *payloads) get(id GroupID, peerId PeerId) Payload {
	payload, _ := p.payloads[id][peerId]
	return payload
}

func (p *payloads) set(id GroupID, peerId PeerId, payload Payload) {
	_, ok := p.payloads[id]
	if !ok {
		p.payloads[id] = make(map[PeerId]Payload)
	}

	p.payloads[id][peerId] = payload
}
