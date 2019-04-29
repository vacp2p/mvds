package mvds

import (
	"github.com/status-im/mvds/securetransport"
	"github.com/status-im/mvds/storage"
)

type calculateSendTime func(count int64, lastTime int64) uint64
type SyncState map[MessageID]State

type State struct {
	HoldFlag    bool
	AckFlag     bool
	RequestFlag bool
	SendCount   int64
	SendTime    uint64
}

type Node struct {
	ms storage.MessageStore
	st securetransport.Node

	ss State

	id []byte // @todo

	Send     <-chan []byte
	Received chan<- []byte // @todo message type
}

func (n *Node) Start() error {

	// @todo start listening to both the send channel and what the transport receives for later handling

	return nil
}

func (n *Node) onPayload(payload Payload) {
	// @todo probably needs to check that its not null and all that
	n.onAck(payload.ack)
	n.onRequest(payload.request)
	n.onOffer(payload.offer)

	for _, m := range payload.messages {
		n.onMessage(m)
	}
}

func (n *Node) onOffer(msg Offer) {

}

func (n *Node) onRequest(msg Request) {
	for _, id := range msg.Messages {
		_, err := n.ms.GetMessage(id)
		if err != nil {
			// @todo
		}

		n.send(id)
	}
}

func (n *Node) onAck(msg Ack) {
	for _, id := range msg.Messages {
		// @todo mark acked
	}
}

func (n *Node) onMessage(msg Message) {
	// @todo handle acks for these messages
}

func (n *Node) send(id MessageID) error {
	return nil
}
