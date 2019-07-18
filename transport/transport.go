// Package transport contains transport related logic for MVDS.
package transport

import (
	"github.com/vacp2p/mvds/protobuf"
	"github.com/vacp2p/mvds/state"
)

type Packet struct {
	Group   state.GroupID
	Sender  state.PeerID
	Payload protobuf.Payload
}

// Transport defines an interface allowing for agnostic transport implementations.
type Transport interface {
	Watch() Packet
	Send(group state.GroupID, sender state.PeerID, peer state.PeerID, payload protobuf.Payload) error
}
