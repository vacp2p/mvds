package mvds

import (
	"sync"

	"github.com/status-im/mvds/protobuf"
)

type payloads struct {
	sync.Mutex

	payloads map[GroupID]map[PeerID]protobuf.Payload
}
// @todo check in all the functions below that we aren't duplicating stuff

func newPayloads() payloads {
	return payloads{
		payloads: make(map[GroupID]map[PeerID]protobuf.Payload),
	}
}

func (p *payloads) AddOffers(group GroupID, peer PeerID, offers ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Offer == nil {
		payload.Offer = &protobuf.Offer{Id: make([][]byte, 0)}
	}

	payload.Offer.Id = append(payload.Offer.Id, offers...)

	p.set(group, peer, payload)
}

func (p *payloads) AddAcks(group GroupID, peer PeerID, acks ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Ack == nil {
		payload.Ack = &protobuf.Ack{Id: make([][]byte, 0)}
	}

	payload.Ack.Id = append(payload.Ack.Id, acks...)

	p.set(group, peer, payload)
}

func (p *payloads) AddRequests(group GroupID, peer PeerID, request ...[]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Request == nil {
		payload.Request = &protobuf.Request{Id: make([][]byte, 0)}
	}

	payload.Request.Id = append(payload.Request.Id, request...)

	p.set(group, peer, payload)
}

func (p *payloads) AddMessages(group GroupID, peer PeerID, messages ...*protobuf.Message) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Messages == nil {
		payload.Messages = make([]*protobuf.Message, 0)
	}

	payload.Messages = append(payload.Messages, messages...)
	p.set(group, peer, payload)
}

func (p *payloads) MapAndClear(f func(GroupID, PeerID, protobuf.Payload)) {
	p.Lock()
	defer p.Unlock()

	for g, payloads := range p.payloads {
		for peer, payload := range payloads {
			f(g, peer, payload)
		}
	}

	p.payloads = make(map[GroupID]map[PeerID]protobuf.Payload)
}

func (p *payloads) get(id GroupID, peer PeerID) protobuf.Payload {
	payload, _ := p.payloads[id][peer]
	return payload
}

func (p *payloads) set(id GroupID, peer PeerID, payload protobuf.Payload) {
	_, ok := p.payloads[id]
	if !ok {
		p.payloads[id] = make(map[PeerID]protobuf.Payload)
	}

	p.payloads[id][peer] = payload
}
