package mvds

type Transport interface {
	Watch() (PeerId, Payload) // @todo might need be changed in the future
	Send(sender PeerId, peer PeerId, payload Payload) error
}
