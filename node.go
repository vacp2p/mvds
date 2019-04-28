package mvds

import (
	"github.com/status-im/mvds/securetransport"
	"github.com/status-im/mvds/storage"
)

type Node struct {
	ms storage.MessageStore
	st securetransport.Node

	id []byte // @todo

	Send     <-chan []byte
	Received chan<- []byte // @todo message type
}

func (n *Node) Start() error {

	// @todo start listening to both the send channel and what the transport receives for later handling

	return nil
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

}


func (n *Node) send(id MessageID) error {
	return nil
}