package mvds

// @todo this is a very rough implementation that needs cleanup

import (
	"time"

	"github.com/status-im/mvds/securetransport"
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
	ms MessageStore
	st securetransport.Node

	syncState       map[MessageID]map[PeerId]*State
	offeredMessages map[PeerId][]MessageID
	sharing         map[GroupID][]PeerId

	sc calculateSendTime

	id    PeerId
	group GroupID
}

func NewNode(ms MessageStore, st securetransport.Node, sc calculateSendTime, id PeerId, group GroupID) Node {
	return Node{
		ms:              ms,
		st:              st,
		syncState:       make(map[MessageID]map[PeerId]*State),
		offeredMessages: make(map[PeerId][]MessageID),
		sharing:         make(map[GroupID][]PeerId),
		sc:              sc,
		id:              id,
		group:           group,
	}
}

func (n *Node) Run() error {

	// @todo start listening to both the send channel and what the transport receives for later handling

	return nil
}

// @todo
func (n *Node) Send(data []byte) error {

	// @todo this will call a function similar to Chat.send();

	return nil
}

func (n *Node) sendMessages() {

	//pls := n.payloads()

	//for id, p := range pls {
		//err := n.st.SendPayload(n.id, id, *p)
		//if err != nil {
		//	@todo
		//}
	//}
}

func (n *Node) onPayload(sender PeerId, payload Payload) {
	// @todo probably needs to check that its not null and all that
	// @todo do these need to be go routines?
	n.onAck(sender, payload.Ack)
	n.onRequest(sender, payload.Request)
	n.onOffer(sender, payload.Offer)

	for _, m := range payload.Messages {
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

	// @todo push message somewhere for end user
}

func (n *Node) payloads() map[PeerId]*Payload {
	pls := make(map[PeerId]*Payload)

	// Ack offered Messages
	for peer, messages := range n.offeredMessages {
		for _, id := range messages {

			// Ack offered Messages
			if n.ms.HasMessage(id) && n.syncState[id][peer].AckFlag {
				n.syncState[id][peer].AckFlag = false
				pls[peer].Ack.Messages = append(pls[peer].Ack.Messages, id)
			}

			// Request offered Messages
			if !n.ms.HasMessage(id) && n.syncState[id][peer].SendTime <= time.Now().Unix() {
				pls[peer].Request.Messages = append(pls[peer].Request.Messages, id)
				n.syncState[id][peer].HoldFlag = true
				n.updateSendTime(id, peer)
			}
		}
	}

	for id, peers := range n.syncState {
		for peer, s := range peers {
			// Ack sent Messages
			if s.AckFlag {
				pls[peer].Ack.Messages = append(pls[peer].Ack.Messages, id)
				s.AckFlag = false
			}

			if n.isPeerInGroup(n.group, peer) && s.SendTime <= time.Now().Unix() {
				// Offer Messages
				if !s.HoldFlag {
					pls[peer].Offer.Messages = append(pls[peer].Offer.Messages, id)
					n.updateSendTime(id, peer)
				}

				// send requested Messages
				if s.RequestFlag {
					m, err := n.ms.GetMessage(id)
					if err != nil {
						// @todo
					}

					pls[peer].Messages = append(pls[peer].Messages, m)
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
