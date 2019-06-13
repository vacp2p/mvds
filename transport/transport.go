package transport

import (
	"github.com/status-im/mvds/protobuf"
	"github.com/status-im/mvds/state"
)

type Packet struct {
	Group   state.GroupID
	Sender  state.PeerID
	Payload protobuf.Payload
}

type Transport interface {
	Watch() Packet
	Send(group state.GroupID, sender state.PeerID, peer state.PeerID, payload protobuf.Payload) error
}
