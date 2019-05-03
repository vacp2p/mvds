package main

import (
	"errors"

	"github.com/status-im/mvds"
)

type packet struct {
	sender  mvds.PeerId
	payload mvds.Payload
}

type Transport struct {
	in  <-chan packet
	out map[mvds.PeerId]chan<- packet
}

func (t *Transport) Watch() (mvds.PeerId, mvds.Payload) {
	p := <-t.in
	return p.sender, p.payload
}

func (t *Transport) Send(sender mvds.PeerId, peer mvds.PeerId, payload mvds.Payload) error {
	c, ok := t.out[peer]
	if !ok {
		return errors.New("peer unknown")
	}

	c <- packet{sender: sender, payload: payload}
	return nil
}
