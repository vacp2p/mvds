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
		if _, ok := n.syncState[id]; !ok || n.syncState[id][sender].AckFlag {
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
	n.syncState[msg.ID()][sender].HoldFlag = true
	n.syncState[msg.ID()][sender].AckFlag = true

	err := n.ms.SaveMessage(msg)
	if err != nil {
		// @todo process, should this function ever even have an error?
	}
}

func (n *Node) payloads() map[PeerId]*Payload {
	pls := make(map[PeerId]*Payload)

	// ack offered messages
	for peer, messages := range n.offeredMessages {
		for _, id := range messages {

			// ack offered messages
			if n.ms.HasMessage(id) && n.syncState[id][peer].AckFlag {
				n.syncState[id][peer].AckFlag = false
				pls[peer].ack.Messages = append(pls[peer].ack.Messages, id)
			}

			// request offered messages
			if !n.ms.HasMessage(id) && n.syncState[id][peer].SendTime <= time.Now().Unix() {
				pls[peer].request.Messages = append(pls[peer].request.Messages, id)
				n.syncState[id][peer].HoldFlag = true
				n.updateSendTime(id, peer)
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
				if !s.HoldFlag {
					pls[peer].offer.Messages = append(pls[peer].offer.Messages, id)
					n.updateSendTime(id, peer)
				}

				// send requested messages
				if s.RequestFlag {
					m, err := n.ms.GetMessage(id)
					if err != nil {
						// @todo
					}

					pls[peer].messages = append(pls[peer].messages, m)
					n.updateSendTime(id, peer)
					s.RequestFlag = false
				}
			}
		}
	}

	return pls
}

func (n *Node) updateSendTime(m MessageID, p PeerId) {
	n.syncState[m][p].SendCount += 1
	n.syncState[m][p].SendTime = n.sc(n.syncState[m][p].SendCount, n.syncState[m][p].SendTime)
}

func (n Node) isPeerInGroup(g GroupID, p PeerId) bool {
	for _, peer := range n.sharing[g] {
		if peer == p {
			return true
		}
	}

	return false
}
