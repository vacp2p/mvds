package mvds

type Transport interface {
	Send(sender PeerId, peer PeerId, payload Payload) error
}
