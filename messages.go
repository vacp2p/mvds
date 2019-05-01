package mvds

import (
	"crypto/sha256"
	"encoding/binary"
)

// @todo: This will probably be changed to protocol buffers

type MessageID [32]byte
type GroupID [32]byte

type Payload struct {
	Ack      Ack
	Offer    Offer
	Request  Request
	Messages []Message
}

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
	GroupID   GroupID
	Timestamp uint64
	Body      []byte
}

func (m Message) ID() MessageID {
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(m.Timestamp))

	b := append([]byte("MESSAGE_ID"), m.GroupID[:]...)
	b = append(b, t...)
	b = append(b, m.Body...)

	return sha256.Sum256(b)
}
