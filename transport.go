package mvds

type Packet struct {
	Group   GroupID
	Sender  PeerID
	Payload Payload
}

type Transport interface {
	Watch() Packet // @todo might need be changed in the future
	Send(group GroupID, sender PeerID, peer PeerID, payload Payload) error
}
