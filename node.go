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
	offeredMessages map[GroupID]map[PeerId][]MessageID
	sharing         map[GroupID][]PeerId
	peers           map[GroupID][]PeerId

	nextEpoch calculateNextEpoch

	ID PeerId

	epoch int64
}

func NewNode(ms MessageStore, st Transport, nextEpoch calculateNextEpoch, id PeerId) *Node {
	return &Node{
		ms:              ms,
		st:              st,
		s:               syncState{state: make(map[GroupID]map[MessageID]map[PeerId]state)},
		offeredMessages: make(map[GroupID]map[PeerId][]MessageID),
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


	for {
		<-time.After(1 * time.Second)

		// @todo should probably do a select statement
		// @todo this is done very badly
		go func() {
			p := n.st.Watch()
			if p == nil {
				return
			}

			n.onPayload(p.Group, p.Sender, p.Payload)
		}()

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

		s := n.s.Get(group, id, sender)
		s.HoldFlag = true
		n.s.Set(group, id, sender, s)

		log.Printf("[%x] OFFER (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onRequest(group GroupID, sender PeerId, msg Request) {
	for _, raw := range msg.Id {
		id := toMessageID(raw)

		s := n.s.Get(group, id, sender)
		s.RequestFlag = true
		n.s.Set(group, id, sender, s)

		log.Printf("[%x] REQUEST (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onAck(group GroupID, sender PeerId, msg Ack) {
	for _, raw := range msg.Id {
		id := toMessageID(raw)

		s := n.s.Get(group, id, sender)
		s.HoldFlag = true
		n.s.Set(group, id, sender, s)

		log.Printf("[%x] ACK (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])
	}
}

func (n *Node) onMessage(group GroupID, sender PeerId, msg Message) {
	id := msg.ID()

	s := n.s.Get(group, id, sender)
	s.HoldFlag = true
	s.AckFlag = true
	n.s.Set(group, id, sender, s)

	err := n.ms.SaveMessage(msg)
	if err != nil {
		// @todo process, should this function ever even have an error?
	}

	log.Printf("[%x] MESSAGE (%x -> %x): %x received.\n", group[:4], sender.toBytes()[:4], n.ID.toBytes()[:4], id[:4])

	// @todo push message somewhere for end user
}

func (n *Node) payloads() map[GroupID]map[PeerId]*Payload {
	n.Lock()
	defer n.Unlock()

	pls := make(map[GroupID]map[PeerId]*Payload)

	// Ack offered Messages
	o := n.offeredMessages
	for group, offers := range o {
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
				if n.ms.HasMessage(id) && n.s.Get(group, id, peer).AckFlag {

					s := n.s.Get(group, id, peer)
					s.AckFlag = true
					n.s.Set(group, id, peer, s)

					pls[group][peer].Ack.Id = append(pls[group][peer].Ack.Id, id[:])
				}

				// Request offered Messages
				if !n.ms.HasMessage(id) && n.s.Get(group, id, peer).SendEpoch <= n.epoch {
					pls[group][peer].Request.Id = append(pls[group][peer].Request.Id, id[:])

					s := n.s.Get(group, id, peer)
					s.HoldFlag = true
					n.s.Set(group, id, peer, s)

					n.updateSendEpoch(group, id, peer)
				}
			}
		}
	}

	for group, syncstate := range n.s.Iterate() {
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
					n.s.Set(group, id, peer, s)
				}

				if n.IsPeerInGroup(group, peer) && s.SendEpoch <= n.epoch {
					// Offer Messages
					if !s.HoldFlag {
						pls[group][peer].Offer.Id = append(pls[group][peer].Offer.Id, id[:])
						n.updateSendEpoch(group, id, peer)

						// @todo do we wanna send messages like in interactive mode?
					}

					// send requested Messages
					if s.RequestFlag {
						m, err := n.ms.GetMessage(id)
						if err != nil {
							// @todo
						}

						pls[group][peer].Messages = append(pls[group][peer].Messages, &m)
						n.updateSendEpoch(group, id, peer)
						s.RequestFlag = false
						n.s.Set(group, id, peer, s)
					}
				}
			}
		}
	}

	return pls
}

func (n *Node) offerMessage(group GroupID, sender PeerId, id MessageID) {
	n.Lock()
	defer n.Unlock()

	if _, ok := n.offeredMessages[group]; !ok {
		n.offeredMessages[group] = make(map[PeerId][]MessageID)
	}

	if _, ok := n.offeredMessages[group][sender]; !ok {
		n.offeredMessages[group][sender] = make([]MessageID, 0)
	}

	n.offeredMessages[group][sender] = append(n.offeredMessages[group][sender], id)
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
