package mvds

// @todo: This will probably be changed to protocol buffers

type MessageID [32]byte

type Ack struct {
	Messages []MessageID
}

type Offer struct {
	Messages []MessageID
}

type Request struct {
	Messages []MessageID
}

type Message struct {
	GroupID   [32]byte
	Timestamp int64
	Body      []byte
}

func (m Message) ID() MessageID {
	// @todo
}

type Payload struct {
	ack      Ack
	offer    Offer
	request  Request
	messages []Message
}
