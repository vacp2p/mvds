package mvds

import "github.com/status-im/mvds/storage"

type Node struct {
	ms storage.MessageStore

	Send     <-chan []byte
	Received chan<- []byte // @todo message type
}

func (n *Node) Start() error {

	// @todo start listening to both the send channel and what the transport receives for later handling

	return nil
}

func (n *Node) onRequest(message Request) {
	for _, id := range message.Messages {
		m, err := n.ms.GetMessage(id)
		if err != nil {
			// @todo
		}

		// @todo put message on the wire
	}
}

func (n *Node) onAck() {

}

func (n *Node) onMessage() {

}
