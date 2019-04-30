package mvds

import (
	"github.com/status-im/mvds/securetransport"
	"github.com/status-im/mvds/storage"
)

type calculateSendTime func(count uint64, lastTime uint64) uint64
type PeerId [32]byte

type State struct {
	HoldFlag    bool
	AckFlag     bool
	RequestFlag bool
	SendCount   uint64
	SendTime    uint64
}

type Node struct {
	ms storage.MessageStore
	st securetransport.Node

	syncState       map[MessageID]map[PeerId]State
	offeredMessages map[PeerId][]MessageID

	queue map[PeerId]Payload // @todo we use this so we can queue messages rather than sending stuff alone
							 // @todo make this a new object which is mutexed

	sc calculateSendTime

	id []byte // @todo

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

		s, _ := n.syncState[id][sender]
		s.HoldFlag = true
	}
}

func (n *Node) onRequest(sender PeerId, msg Request) {
	for _, id := range msg.Messages {
		s, _ := n.syncState[id][sender]
		s.RequestFlag = true
	}
}

func (n *Node) onAck(sender PeerId, msg Ack) {
	for _, id := range msg.Messages {
		s, _ := n.syncState[id][sender]
		s.HoldFlag = true
	}
}

func (n *Node) onMessage(sender PeerId, msg Message) {

	// @todo handle

	n.syncState[msg.ID()][sender] = State{
		HoldFlag:    true,
		AckFlag:     true,
		RequestFlag: false,
		SendTime:    0,
		SendCount:   0,
	}

	err := n.ms.SaveMessage(msg)
	if err != nil {
		// @todo process, should this function ever even have an error?
	}
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
