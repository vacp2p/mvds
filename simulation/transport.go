package simulation

import (
	"errors"

	"github.com/status-im/mvds"
)

type Transport struct {
	out map[mvds.PeerId]chan mvds.Payload
}

func (t *Transport) Send(sender mvds.PeerId, peer mvds.PeerId, payload mvds.Payload) error {
	c, ok := t.out[peer]
	if !ok {
		return errors.New("peer unknown")
	}

	c <- payload
	return nil
}

