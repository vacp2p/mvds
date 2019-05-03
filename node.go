package mvds

// @todo this is a very rough implementation that needs cleanup

import (
	"time"
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
	st Transport

	syncState       map[MessageID]map[PeerId]*State
	offeredMessages map[PeerId][]MessageID
	sharing         map[GroupID][]PeerId
	peers           []PeerId

	sc calculateSendTime

	ID    PeerId
	group GroupID

	time int64
}

func NewNode(ms MessageStore, st Transport, sc calculateSendTime, id PeerId, group GroupID) Node {
	return Node{
		ms:              ms,
		st:              st,
		syncState:       make(map[MessageID]map[PeerId]*State),
		offeredMessages: make(map[PeerId][]MessageID),
		sharing:         make(map[GroupID][]PeerId),
		peers:           make([]PeerId, 0),
		sc:              sc,
		ID:              id,
		group:           group,
		time:            0,
	}
}

func (n *Node) Run() {

	// @todo start listening to both the send channel and what the transport receives for later handling

	// @todo maybe some waiting?
	for {
		s, p := n.st.Watch()
		go n.onPayload(s, p)

		go n.sendMessages() // @todo probably not that efficient here

		n.time += 1
	}
}

func (n *Node) Send(data []byte) error {
	m := Message{
		GroupId:   n.group[:],
		Timestamp: time.Now().Unix(),
		Body:      data,
	}

	err := n.ms.SaveMessage(m)
	if err != nil {
		return err
	}

	id := m.ID()

	for _, p := range n.peers {
		if !n.isPeerInGroup(n.group, p) {
			continue
		}

		s := n.state(id, p)
		s.SendTime = n.time + 1
	}

	// @todo think about a way to insta trigger send messages when send was selected, we don't wanna wait for ticks here

	return nil
}

func (n *Node) sendMessages() {

	pls := n.payloads()

	for id, p := range pls {
		err := n.st.Send(n.ID, id, *p)
		if err != nil {
			//	@todo
		}
	}
}

func (n *Node) onPayload(sender PeerId, payload Payload) {
	// @todo probably needs to check that its not null and all that
	// @todo do these need to be go routines?
	n.onAck(sender, *payload.Ack)
	n.onRequest(sender, *payload.Request)
	n.onOffer(sender, *payload.Offer)

	for _, m := range payload.Messages {
		n.onMessage(sender, *m)
	}
}

func (n *Node) onOffer(sender PeerId, msg Offer) {
	for _, raw := range msg.Id {
		id := toMessageID(raw)

		if _, ok := n.syncState[id]; !ok || n.syncState[id][sender].AckFlag {
			n.appendOfferedMessage(sender, id)
		}

		n.state(id, sender).HoldFlag = true
	}
}

func (n *Node) onRequest(sender PeerId, msg Request) {
	for _, id := range msg.Id {
		n.state(toMessageID(id), sender).RequestFlag = true
	}
}

func (n *Node) onAck(sender PeerId, msg Ack) {
	for _, id := range msg.Id {
		n.state(toMessageID(id), sender).HoldFlag = true
	}
}

func (n *Node) onMessage(sender PeerId, msg Message) {
	s := n.state(msg.ID(), sender)
	s.HoldFlag = true
	s.AckFlag = true

	err := n.ms.SaveMessage(msg)
	if err != nil {
		// @todo process, should this function ever even have an error?
	}

	// @todo push message somewhere for end user
}

func (n *Node) payloads() map[PeerId]*Payload {

	// @todo do we need a mutex?

	pls := make(map[PeerId]*Payload)

	// Ack offered Messages
	for peer, messages := range n.offeredMessages {
		for _, id := range messages {

			// Ack offered Messages
			if n.ms.HasMessage(id) && n.syncState[id][peer].AckFlag {
				n.syncState[id][peer].AckFlag = false
				pls[peer].Ack.Id = append(pls[peer].Ack.Id, id[:])
			}

			// Request offered Messages
			if !n.ms.HasMessage(id) && n.syncState[id][peer].SendTime <= n.time {
				pls[peer].Request.Id = append(pls[peer].Request.Id, id[:])
				n.syncState[id][peer].HoldFlag = true
				n.updateSendTime(id, peer)
			}
		}
	}

	for id, peers := range n.syncState {
		for peer, s := range peers {
			// Ack sent Messages
			if s.AckFlag {
				pls[peer].Ack.Id = append(pls[peer].Ack.Id, id[:])
				s.AckFlag = false
			}

			if n.isPeerInGroup(n.group, peer) && s.SendTime <= n.time {
				// Offer Messages
				if !s.HoldFlag {
					pls[peer].Offer.Id = append(pls[peer].Offer.Id, id[:])
					n.updateSendTime(id, peer)
				}

				// send requested Messages
				if s.RequestFlag {
					m, err := n.ms.GetMessage(id)
					if err != nil {
						// @todo
					}

					pls[peer].Messages = append(pls[peer].Messages, &m)
					n.updateSendTime(id, peer)
					s.RequestFlag = false
				}
			}
		}
	}

	return pls
}

func (n *Node) state(id MessageID, sender PeerId) *State {
	if _, ok := n.syncState[id]; !ok {
		n.syncState[id] = make(map[PeerId]*State)
	}

	if _, ok := n.syncState[id][sender]; !ok {
		n.syncState[id][sender] = &State{}
	}

	return n.syncState[id][sender]
}

func (n *Node) appendOfferedMessage(sender PeerId, id MessageID) {
	if _, ok := n.offeredMessages[sender]; !ok {
		n.offeredMessages[sender] = make([]MessageID, 0)
	}

	n.offeredMessages[sender] = append(n.offeredMessages[sender], id)
}

func (n *Node) updateSendTime(m MessageID, p PeerId) {
	s := n.state(m, p)
	s.SendCount += 1
	s.SendTime = n.sc(s.SendCount, s.SendTime)
}

func (n Node) isPeerInGroup(g GroupID, p PeerId) bool {
	for _, peer := range n.sharing[g] {
		if peer == p {
			return true
		}
	}

	return false
}

func toMessageID(b []byte) MessageID {
	var id MessageID
	copy(id[:], b)
	return id
}
