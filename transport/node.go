package transport

type Node interface {
	Watch()
	SendMessage(senderId []byte, to []byte, message []byte) error // @todo probably needs types
}
