package transport

type Node interface {
	Watch() error // @todo this should return error or payload, watch should always watch for one message and then be called again
	SendMessage(senderId []byte, to []byte, message []byte) error // @todo probably needs types
}
