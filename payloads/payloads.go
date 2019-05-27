package payloads

import (
	"sync"

	"github.com/status-im/mvds"
)

type Payloads struct {
	sync.RWMutex

	payloads map[mvds.GroupID]map[mvds.PeerId]mvds.Payload
}

func (p *Payloads) get(id mvds.GroupID, peerId mvds.PeerId) mvds.Payload {
	payload, _ := p.payloads[id][peerId]
	return payload
}

func (p *Payloads) set(id mvds.GroupID, peerId mvds.PeerId, payload mvds.Payload) {
	p.payloads[id][peerId] = payload
}

func (p *Payloads) AddOffers(group mvds.GroupID, peer mvds.PeerId, offers [][]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Offer == nil {
		payload.Offer = &mvds.Offer{Id: make([][]byte, 0)}
	}

	payload.Offer.Id = append(payload.Offer.Id, offers...)
	p.set(group, peer, payload)
}

func (p *Payloads) AddAcks(group mvds.GroupID, peer mvds.PeerId, acks [][]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Ack == nil {
		payload.Ack = &mvds.Ack{Id: make([][]byte, 0)}
	}

	payload.Ack.Id = append(payload.Ack.Id, acks...)
	p.set(group, peer, payload)
}

func (p *Payloads) AddRequests(group mvds.GroupID, peer mvds.PeerId, request [][]byte) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Request == nil {
		payload.Request = &mvds.Request{Id: make([][]byte, 0)}
	}

	payload.Request.Id = append(payload.Request.Id, request...)
	p.set(group, peer, payload)
}

func (p *Payloads) AddMessages(group mvds.GroupID, peer mvds.PeerId, messages []*mvds.Message) {
	p.Lock()
	defer p.Unlock()

	payload := p.get(group, peer)
	if payload.Messages == nil {
		payload.Messages = make([]*mvds.Message, 0)
	}

	payload.Messages = append(payload.Messages, messages...)
	p.set(group, peer, payload)
}

func (p *Payloads) Map(f func(mvds.GroupID, mvds.PeerId, mvds.Payload)) {
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

	p.payloads = make(map[mvds.GroupID]map[mvds.PeerId]mvds.Payload)
}
