package transport

type Node interface {
	Watch()
	SendMessage(senderId []byte, to []byte, message []byte) // @todo probably needs types
}
