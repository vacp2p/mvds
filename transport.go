package mvds

type Packet struct {
	Group   GroupID
	Sender  PeerId
	Payload Payload
}

type Transport interface {
	Watch() Packet // @todo might need be changed in the future
	Send(sender PeerId, peer PeerId, payload Payload) error
}
