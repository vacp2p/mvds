package transport

type Node interface {
	Tick()
	SendMessage(senderId []byte, to []byte, message []byte) // @todo probably needs types
}
