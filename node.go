package mvds

// @todo this is a very rough implementation that needs cleanup

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type calculateNextEpoch func(count uint64, epoch int64) int64
type PeerId ecdsa.PublicKey

type State struct {
	HoldFlag    bool
	AckFlag     bool
	RequestFlag bool
	SendCount   uint64
	SendEpoch   int64
}

type Node struct {
	sync.Mutex

	ms MessageStore
	st Transport

	s               syncState
	sharing         map[GroupID][]PeerId
	peers           map[GroupID][]PeerId

	payloads map[GroupID]map[PeerId]*Payload

	nextEpoch calculateNextEpoch

	ID PeerId

	epoch int64
}

func NewNode(ms MessageStore, st Transport, nextEpoch calculateNextEpoch, id PeerId) *Node {
	return &Node{
		ms:              ms,
		st:              st,
		s:               syncState{state: make(map[GroupID]map[MessageID]map[PeerId]state)},
		sharing:         make(map[GroupID][]PeerId),
		peers:           make(map[GroupID][]PeerId),
		nextEpoch:       nextEpoch,
		ID:              id,
		epoch:           0,
	}
}

// Run listens for new messages received by the node and sends out those required every tick.
func (n *Node) Run() {

	// @todo start listening to both the send channel and what the transport receives for later handling

	// @todo maybe some waiting?

	// this will be completely legitimate with new payload handling
	go func() {
		for {
			p := n.st.Watch()
			n.onPayload(p.Group, p.Sender, p.Payload)
		}
	}()

	for {
		<-time.After(1 * time.Second)

		go n.sendMessages() // @todo probably not that efficient here
		n.epoch += 1
	}
}

// AppendMessage sends a message to a given group.
func (n *Node) AppendMessage(group GroupID, data []byte) (MessageID, error) {

	// @todo because we don't lock here we seem to get to a point where we can no longer append?

	m := Message{
		GroupId:   group[:],
		Timestamp: time.Now().Unix(),
		Body:      data,
	}

	err := n.ms.SaveMessage(m)
	if err != nil {
		return MessageID{}, err
	}

	id := m.ID()

	for g, peers := range n.peers {
		for _, p := range peers {
			if !n.IsPeerInGroup(group, p) {
				continue
			}

			s := n.s.Get(g, id, p)
			s.SendEpoch = n.epoch + 1
			n.s.Set(g, id, p, s)
		}
	}

	// @todo think about a way to insta trigger send messages when send was selected, we don't wanna wait for ticks here

	return id, nil
}

// AddPeer adds a peer to a specific group making it a recipient of messages
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
	m := n.onRequest(group, sender, *payload.Request)
	r := n.onOffer(group, sender, *payload.Offer)
	a := n.onMessages(group, sender, payload.Messages)

	n.payloads[group][sender] = &Payload{
		Ack: &a,
		Offer: &Offer{Id: make([][]byte, 0)},
		Request: &r,
		Messages: m,
	}
}

func (n *Node) onOffer(group GroupID, sender PeerId, msg Offer) Request {
	r := Request{Id: make([][]byte, 0)}

	for _, raw := range msg.Id {
		id := toMessageID(raw)
		log.Printf("[%x] OFFER (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

		// @todo maybe ack?
		if n.ms.HasMessage(id) {
			continue
		}

		r.Id = append(r.Id, raw)
		log.Printf("[%x] sending REQUEST (%x -> %x): %x\n", group[:4], n.ID.toBytes()[:4], sender.toBytes()[:4], id[:4])
	}

	return r
}

func (n *Node) onRequest(group GroupID, sender PeerId, msg Request) []*Message {
	m := make([]*Message, 0)

	for _, raw := range msg.Id {
		id := toMessageID(raw)
		log.Printf("[%x] REQUEST (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

		message, err := n.ms.GetMessage(id)
		if err != nil {
			log.Printf("error requesting message %x", id[:4])
			continue
		}

		// @todo should ensure we are sharing, we don't wanna just send anyone anything

		// @todo send count and send epoch
		// s.SendCount += 1
		// s.SendEpoch += n.nextEpoch(s.SendCount, n.epoch)

		m = append(m, &message)
	}

	return m
}

// @todo this should return nothing?
func (n *Node) onAck(group GroupID, sender PeerId, msg Ack) {
	for _, raw := range msg.Id {
		id := toMessageID(raw)

		s := n.s.Get(group, id, sender)
		s.HoldFlag = true
		n.s.Set(group, id, sender, s)

		log.Printf("[%x] ACK (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onMessages(group GroupID, sender PeerId, messages []*Message) Ack {
	a := Ack{Id: make([][]byte, 0)}

	for _, m := range messages {
		err := n.onMessage(group, sender, *m)
		if err != nil {
			// @todo
			continue
		}

		id := m.ID()
		log.Printf("[%x] sending ACK (%x -> %x): %x\n", group[:4], n.ID.toBytes()[:4], sender.toBytes()[:4], id[:4])
		a.Id = append(a.Id, id[:])
	}

	return a
}

// @todo this should return ACKs
func (n *Node) onMessage(group GroupID, sender PeerId, msg Message) error {
	id := msg.ID()
	log.Printf("[%x] MESSAGE (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

	s := n.s.Get(group, id, sender)
	s.HoldFlag = true
	s.AckFlag = true
	n.s.Set(group, id, sender, s)

	err := n.ms.SaveMessage(msg)
	if err != nil {
		return err
		// @todo process, should this function ever even have an error?
	}


	// @todo push message somewhere for end user
	return nil
}

func (n *Node) updateSendEpoch(g GroupID, m MessageID, p PeerId) {
	s := n.s.Get(g, m, p)
	s.SendCount += 1
	s.SendEpoch += n.nextEpoch(s.SendCount, n.epoch)
	n.s.Set(g, m, p, s)
}

func (n Node) IsPeerInGroup(g GroupID, p PeerId) bool {
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

func (p PeerId) toBytes() []byte {
	if p.X == nil || p.Y == nil {
		return nil
	}

	return elliptic.Marshal(crypto.S256(), p.X, p.Y)
}
