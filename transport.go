package mvds

import "github.com/status-im/mvds/protobuf"

type Packet struct {
	Group   GroupID
	Sender  PeerID
	Payload protobuf.Payload
}

type Transport interface {
	Watch() Packet // @todo might need be changed in the future
	Send(group GroupID, sender PeerID, peer PeerID, payload protobuf.Payload) error
}
