package mvds

// @todo this is a very rough implementation that needs cleanup

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"log"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type calculateNextEpoch func(count uint64, epoch int64) int64
type PeerId ecdsa.PublicKey

type Node struct {
	store     MessageStore
	transport Transport

	syncState syncState
	sharing   map[GroupID][]PeerId
	peers     map[GroupID][]PeerId

	payloads payloads

	nextEpoch calculateNextEpoch

	ID PeerId

	epoch int64
}

func NewNode(ms MessageStore, st Transport, nextEpoch calculateNextEpoch, id PeerId) *Node {
	return &Node{
		store:     ms,
		transport: st,
		syncState: newSyncState(),
		sharing:   make(map[GroupID][]PeerId),
		peers:     make(map[GroupID][]PeerId),
		payloads:  payloads{payloads: make(map[GroupID]map[PeerId]Payload)},
		nextEpoch: nextEpoch,
		ID:        id,
		epoch:     0,
	}
}

// Run listens for new messages received by the node and sends out those required every tick.
func (n *Node) Run() {

	// this will be completely legitimate with new payload handling
	go func() {
		for {
			p := n.transport.Watch()
			go n.onPayload(p.Group, p.Sender, p.Payload)
		}
	}()

	go func() {
		for {
			log.Printf("Node: %x Epoch: %d", n.ID.toBytes()[:4], n.epoch)
			time.Sleep(1 * time.Second)

			n.sendMessages()
			atomic.AddInt64(&n.epoch, 1)
		}
	}()
}

// AppendMessage sends a message to a given group.
func (n *Node) AppendMessage(group GroupID, data []byte) (MessageID, error) {

	// @todo because we don't lock here we seem to get to a point where we can no longer append?

	m := Message{
		GroupId:   group[:],
		Timestamp: time.Now().Unix(),
		Body:      data,
	}

	err := n.store.Add(m)
	if err != nil {
		return MessageID{}, err
	}

	id := m.ID()

	go func () {
		for g, peers := range n.peers {
			for _, p := range peers {
				if !n.IsPeerInGroup(group, p) {
					continue
				}

				// @todo store a sync state only for Offers

				s := state{}
				s.SendEpoch = n.epoch + 1
				n.syncState.Set(g, id, p, s)
			}
		}
	}()

	log.Printf("[%x] node %x sending %x\n", group[:4], n.ID.toBytes()[:4], id[:4])
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

func (n Node) IsPeerInGroup(g GroupID, p PeerId) bool {
	for _, peer := range n.sharing[g] {
		if peer == p {
			return true
		}
	}

	return false
}

func (n *Node) sendMessages() {
	n.syncState.Map(func(g GroupID, m MessageID, p PeerId, s state) state {
		if s.SendEpoch < n.epoch || !n.IsPeerInGroup(g, p) {
			return s
		}

		n.payloads.AddOffers(g, p, m[:])
		return n.updateSendEpoch(s)
	})

	n.payloads.Map(func(id GroupID, peer PeerId, payload Payload) {
		err := n.transport.Send(id, n.ID, peer, payload)
		if err != nil {
			//	@todo
		}
	})

	n.payloads.RemoveAll()
}

func (n *Node) onPayload(group GroupID, sender PeerId, payload Payload) {
	if payload.Ack != nil {
		n.onAck(group, sender, *payload.Ack)
	}

	if payload.Request != nil {
		n.payloads.AddMessages(group, sender, n.onRequest(group, sender, *payload.Request)...)
	}

	if payload.Offer != nil {
		n.payloads.AddRequests(group, sender, n.onOffer(group, sender, *payload.Offer)...)
	}

	if payload.Messages != nil {
		n.payloads.AddAcks(group, sender, n.onMessages(group, sender, payload.Messages)...)
	}
}

func (n *Node) onOffer(group GroupID, sender PeerId, msg Offer) [][]byte {
	r := make([][]byte, 0)

	for _, raw := range msg.Id {
		id := toMessageID(raw)
		log.Printf("[%x] OFFER (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

		// @todo maybe ack?
		if n.store.Has(id) {
			continue
		}

		r = append(r, raw)
		log.Printf("[%x] sending REQUEST (%x -> %x): %x\n", group[:4], n.ID.toBytes()[:4], sender.toBytes()[:4], id[:4])
	}

	return r
}

func (n *Node) onRequest(group GroupID, sender PeerId, msg Request) []*Message {
	m := make([]*Message, 0)

	for _, raw := range msg.Id {
		id := toMessageID(raw)
		log.Printf("[%x] REQUEST (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

		if !n.IsPeerInGroup(group, sender) {
			continue
		}

		message, err := n.store.Get(id)
		if err != nil {
			log.Printf("error requesting message %x", id[:4])
			continue
		}

		n.syncState.Set(group, id, sender, n.updateSendEpoch(n.syncState.Get(group, id, sender)))

		m = append(m, &message)

		log.Printf("[%x] sending MESSAGE (%x -> %x): %x\n", group[:4], n.ID.toBytes()[:4], sender.toBytes()[:4], id[:4])
	}

	return m
}

// @todo this should return nothing?
func (n *Node) onAck(group GroupID, sender PeerId, msg Ack) {
	for _, raw := range msg.Id {
		id := toMessageID(raw)

		n.syncState.Remove(group, id, sender)

		log.Printf("[%x] ACK (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onMessages(group GroupID, sender PeerId, messages []*Message) [][]byte {
	a := make([][]byte, 0)

	for _, m := range messages {
		err := n.onMessage(group, sender, *m)
		if err != nil {
			// @todo
			continue
		}

		id := m.ID()
		log.Printf("[%x] sending ACK (%x -> %x): %x\n", group[:4], n.ID.toBytes()[:4], sender.toBytes()[:4], id[:4])
		a = append(a, id[:])
	}

	return a
}

// @todo this should return ACKs
func (n *Node) onMessage(group GroupID, sender PeerId, msg Message) error {
	id := msg.ID()
	log.Printf("[%x] MESSAGE (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

	// @todo share message with those around us

	err := n.store.Add(msg)
	if err != nil {
		return err
		// @todo process, should this function ever even have an error?
	}

	// @todo push message somewhere for end user
	return nil
}

func (n Node) updateSendEpoch(s state) state {
	s.SendCount += 1
	s.SendEpoch += n.nextEpoch(s.SendCount, n.epoch)
	return s
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
