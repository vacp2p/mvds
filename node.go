package mvds

import (
	"time"

	"github.com/status-im/mvds/securetransport"
	"github.com/status-im/mvds/storage"
)

type calculateSendTime func(count uint64, lastTime int64) int64
type PeerId [32]byte

type State struct {
	HoldFlag    bool
	AckFlag     bool
	RequestFlag bool
	SendCount   uint64
	SendTime    int64
}

type Node struct {
	ms storage.MessageStore
	st securetransport.Node

	syncState       map[MessageID]map[PeerId]*State
	offeredMessages map[PeerId][]MessageID
	sharing         map[GroupID][]PeerId

	sc calculateSendTime

	id    PeerId
	group GroupID

	Send     <-chan []byte
	Received chan<- []byte // @todo message type
}

func (n *Node) Start() error {

	// @todo start listening to both the send channel and what the transport receives for later handling

	return nil
}

func (n *Node) onPayload(sender PeerId, payload Payload) {
	// @todo probably needs to check that its not null and all that
	// @todo do these need to be go routines?
	n.onAck(sender, payload.ack)
	n.onRequest(sender, payload.request)
	n.onOffer(sender, payload.offer)

	for _, m := range payload.messages {
		n.onMessage(sender, m)
	}
}

func (n *Node) onOffer(sender PeerId, msg Offer) {
	for _, id := range msg.Messages {
		if _, ok := n.syncState[id]; !ok || n.syncState[id][sender].AckFlag == true {
			n.offeredMessages[sender] = append(n.offeredMessages[sender], id)
		}

		n.syncState[id][sender].HoldFlag = true
	}
}

func (n *Node) onRequest(sender PeerId, msg Request) {
	for _, id := range msg.Messages {
		n.syncState[id][sender].RequestFlag = true
	}
}

func (n *Node) onAck(sender PeerId, msg Ack) {
	for _, id := range msg.Messages {
		n.syncState[id][sender].HoldFlag = true
	}
}

func (n *Node) onMessage(sender PeerId, msg Message) {
	// @todo do we need to initialize stucts?
	n.syncState[msg.ID()][sender].HoldFlag = true
	n.syncState[msg.ID()][sender].AckFlag = true

	// @todo handle for group

	err := n.ms.SaveMessage(msg)
	if err != nil {
		// @todo process, should this function ever even have an error?
	}
}

func (n *Node) payloads() map[PeerId]Payload {
	pls := make(map[PeerId]Payload)

	// ack offered messages
	for peer, messages := range n.offeredMessages {
		for _, id := range messages {
			if n.ms.HasMessage(id) && n.syncState[id][peer].AckFlag == true {
				n.syncState[id][peer].AckFlag = false
				pls[peer].ack.Messages = append(pls[peer].ack.Messages, id)
			}
		}
	}

	for id, peers := range n.syncState {
		for peer, s := range peers {
			// ack sent messages
			if s.AckFlag {
				pls[peer].ack.Messages = append(pls[peer].ack.Messages, id)
				s.AckFlag = false
			}

			if n.isPeerInGroup(n.group, peer) && s.SendTime <= time.Now().Unix() {
				// offer messages
				if s.HoldFlag == false {
					pls[peer].offer.Messages = append(pls[peer].offer.Messages, id)
					n.syncState[id][peer].SendCount += 1
					n.syncState[id][peer].SendTime = n.sc(n.syncState[id][peer].SendCount, n.syncState[id][peer].SendTime)
				}

				// send requested messages
				if s.RequestFlag == true {
					// @todo send requested messages
				}
			}

			// @todo request offered messages
		}
	}

	return pls
}

func (n Node) isPeerInGroup(g GroupID, p PeerId) bool {
	for _, peer := range n.sharing[g] {
		if peer == p {
			return true
		}
	}

	return false
}

func (n *Node) send(to PeerId, id MessageID) error {
	s, _ := n.syncState[id][to]

	s.SendCount += 1
	s.SendTime = n.sc(s.SendCount, s.SendTime)

	// @todo actually send

	return nil
}

func (n *Node) sendForPeer(peer PeerId) {

}
