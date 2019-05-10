package mvds

// @todo this is a very rough implementation that needs cleanup

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type calculateSendTime func(count uint64, time int64) int64
type PeerId ecdsa.PublicKey

type State struct {
	HoldFlag    bool
	AckFlag     bool
	RequestFlag bool
	SendCount   uint64
	SendTime    int64
}

type Node struct {
	sync.Mutex

	ms MessageStore
	st Transport

	syncState       map[GroupID]map[MessageID]map[PeerId]*State
	offeredMessages map[GroupID]map[PeerId][]MessageID
	sharing         map[GroupID][]PeerId
	peers           map[GroupID][]PeerId

	sc calculateSendTime

	ID PeerId

	time int64
}

func NewNode(ms MessageStore, st Transport, sc calculateSendTime, id PeerId) Node {
	return Node{
		ms:              ms,
		st:              st,
		syncState:       make(map[GroupID]map[MessageID]map[PeerId]*State),
		offeredMessages: make(map[GroupID]map[PeerId][]MessageID),
		sharing:         make(map[GroupID][]PeerId),
		peers:           make(map[GroupID][]PeerId),
		sc:              sc,
		ID:              id,
		time:            0,
	}
}

func (n *Node) Run() {

	// @todo start listening to both the send channel and what the transport receives for later handling

	// @todo maybe some waiting?

	for {
		<-time.After(1 * time.Second)

		// @todo should probably do a select statement
		go func() {
			p := n.st.Watch()
			n.onPayload(p.Group, p.Sender, p.Payload)
		}()

		go n.sendMessages() // @todo probably not that efficient here
		n.time += 1
	}
}

func (n *Node) Send(group GroupID, data []byte) error {
	n.Lock()
	defer n.Unlock()

	m := Message{
		GroupId:   group[:],
		Timestamp: time.Now().Unix(),
		Body:      data,
	}

	err := n.ms.SaveMessage(m)
	if err != nil {
		return err
	}

	id := m.ID()

	for g, peers := range n.peers {
		for _, p := range peers {
			if !n.isPeerInGroup(group, p) {
				continue
			}

			n.state(g, id, p).SendTime = n.time + 1
		}
	}

	// @todo think about a way to insta trigger send messages when send was selected, we don't wanna wait for ticks here

	return nil
}

func (n *Node) AddPeer(group GroupID, id PeerId) {
	if _, ok := n.peers[group]; !ok {
		n.peers[group] = make([]PeerId, 0)
	}

	n.peers[group] = append(n.peers[group], id)
}

func (n *Node) Share(group GroupID, id PeerId) {
	if _, ok := n.sharing[group]; !ok {
		n.sharing[group] = make([]PeerId, 0)
	}

	n.sharing[group] = append(n.sharing[group], id)
}

func (n *Node) sendMessages() {

	pls := n.payloads()

	for g, payloads := range pls {
		for id, p := range payloads {

			err := n.st.Send(g, n.ID, id, *p)
			if err != nil {
				//	@todo
			}
		}
	}
}

func (n *Node) onPayload(group GroupID, sender PeerId, payload Payload) {
	n.onAck(group, sender, *payload.Ack)
	n.onRequest(group, sender, *payload.Request)
	n.onOffer(group, sender, *payload.Offer)

	for _, m := range payload.Messages {
		n.onMessage(group, sender, *m)
	}
}

func (n *Node) onOffer(group GroupID, sender PeerId, msg Offer) {
	for _, raw := range msg.Id {
		id := toMessageID(raw)
		n.offerMessage(group, sender, id)
		n.state(group, id, sender).HoldFlag = true
		fmt.Printf("[%s] OFFER (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onRequest(group GroupID, sender PeerId, msg Request) {
	for _, id := range msg.Id {
		n.state(group, toMessageID(id), sender).RequestFlag = true
		fmt.Printf("[%s] REQUEST (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onAck(group GroupID, sender PeerId, msg Ack) {
	for _, id := range msg.Id {
		n.state(group, toMessageID(id), sender).HoldFlag = true
		fmt.Printf("[%s] ACK (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onMessage(group GroupID, sender PeerId, msg Message) {
	id := msg.ID()
	s := n.state(group, id, sender)
	s.HoldFlag = true
	s.AckFlag = true

	err := n.ms.SaveMessage(msg)
	if err != nil {
		// @todo process, should this function ever even have an error?
	}

	fmt.Printf("[%s] MESSAGE (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

	// @todo push message somewhere for end user
}

func (n *Node) payloads() map[GroupID]map[PeerId]*Payload {
	n.Lock()
	defer n.Unlock()

	pls := make(map[GroupID]map[PeerId]*Payload)

	// Ack offered Messages
	for group, offers := range n.offeredMessages {
		for peer, messages := range offers {
			// @todo do we need this?
			if _, ok := pls[group]; !ok {
				pls[group] = make(map[PeerId]*Payload)
			}

			if _, ok := pls[group][peer]; !ok {
				pls[group][peer] = createPayload()
			}

			for _, id := range messages {
				// Ack offered Messages
				if n.ms.HasMessage(id) && n.syncState[group][id][peer].AckFlag {
					n.syncState[group][id][peer].AckFlag = false
					pls[group][peer].Ack.Id = append(pls[group][peer].Ack.Id, id[:])
				}

				// Request offered Messages
				if !n.ms.HasMessage(id) && n.state(group, id, peer).SendTime <= n.time {
					pls[group][peer].Request.Id = append(pls[group][peer].Request.Id, id[:])
					n.syncState[group][id][peer].HoldFlag = true
					n.updateSendTime(group, id, peer)
				}
			}
		}
	}

	for group, syncstate := range n.syncState {
		for id, peers := range syncstate {
			for peer, s := range peers {
				if _, ok := pls[group]; !ok {
					pls[group] = make(map[PeerId]*Payload)
				}

				if _, ok := pls[group][peer]; !ok {
					pls[group][peer] = createPayload()
				}

				// Ack sent Messages
				if s.AckFlag {
					pls[group][peer].Ack.Id = append(pls[group][peer].Ack.Id, id[:])
					s.AckFlag = false
				}

				if n.isPeerInGroup(group, peer) && s.SendTime <= n.time {
					// Offer Messages
					if !s.HoldFlag {
						pls[group][peer].Offer.Id = append(pls[group][peer].Offer.Id, id[:])
						n.updateSendTime(group, id, peer)

						// @todo do we wanna send messages like in interactive mode?
					}

					// send requested Messages
					if s.RequestFlag {
						m, err := n.ms.GetMessage(id)
						if err != nil {
							// @todo
						}

						pls[group][peer].Messages = append(pls[group][peer].Messages, &m)
						n.updateSendTime(group, id, peer)
						s.RequestFlag = false
					}
				}
			}
		}
	}

	return pls
}

func (n *Node) state(group GroupID, id MessageID, sender PeerId) *State {
	//n.Lock()
	//defer n.Unlock()

	// @todo check if we need this
	if _, ok := n.syncState[group]; !ok {
		n.syncState[group] = make(map[MessageID]map[PeerId]*State)
	}

	if _, ok := n.syncState[group][id]; !ok {
		n.syncState[group][id] = make(map[PeerId]*State)
	}

	if _, ok := n.syncState[group][id][sender]; !ok {
		n.syncState[group][id][sender] = &State{}
	}

	return n.syncState[group][id][sender]
}

func (n *Node) offerMessage(group GroupID, sender PeerId, id MessageID) {
	if _, ok := n.offeredMessages[group]; !ok {
		n.offeredMessages[group] = make(map[PeerId][]MessageID)
	}

	if _, ok := n.offeredMessages[group][sender]; !ok {
		n.offeredMessages[group][sender] = make([]MessageID, 0)
	}

	n.offeredMessages[group][sender] = append(n.offeredMessages[group][sender], id)
}

func (n *Node) updateSendTime(g GroupID, m MessageID, p PeerId) {
	s := n.state(g, m, p)
	s.SendCount += 1
	s.SendTime = n.sc(s.SendCount, n.time)
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

func createPayload() *Payload {
	return &Payload{
		Ack:      &Ack{Id: make([][]byte, 0)},
		Offer:    &Offer{Id: make([][]byte, 0)},
		Request:  &Request{Id: make([][]byte, 0)},
		Messages: make([]*Message, 0),
	}
}

func (p PeerId) toBytes() []byte {
	if p.X == nil || p.Y == nil {
		return nil
	}

	return elliptic.Marshal(crypto.S256(), p.X, p.Y)
}

