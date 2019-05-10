package mvds

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"reflect"
	"testing"
	"time"
)

func TestOnRequest(t *testing.T) {
	m := randomMessageId()
	g := GroupID{}

	r := Request{}
	r.Id = append(r.Id, m[:])

	n := getNodeForMessageHandlerTest()

	p := randomPeerId()
	n.onRequest(g, p, r)

	if n.state(g, m, p).RequestFlag != true {
		t.Errorf("did not set Request flag to true")
	}
}

func TestOnAck(t *testing.T) {
	m := randomMessageId()
	g := GroupID{}

	a := Ack{}
	a.Id = append(a.Id, m[:])

	n := getNodeForMessageHandlerTest()

	p := randomPeerId()
	n.onAck(g, p, a)

	if n.state(g, m, p).HoldFlag != true {
		t.Errorf("did not set Hold flag to true")
	}
}

func TestOnOffer(t *testing.T) {
	m := randomMessageId()
	g := GroupID{}

	o := Offer{}
	o.Id = append(o.Id, m[:])

	n := getNodeForMessageHandlerTest()

	p := randomPeerId()
	n.onOffer(g, p, o)

	if n.state(g, m, p).HoldFlag != true {
		t.Errorf("did not set Hold flag to true")
	}

	if n.offeredMessages[g][p][0] != m {
		t.Errorf("message was not added to offered list")
	}
}

func TestOnMessage(t *testing.T) {
	n := getNodeForMessageHandlerTest()
	g := GroupID{}

	ds := NewDummyStore()
	n.ms = &ds

	id := randomMessageId()

	m := Message{
		GroupId: id[:],
		Timestamp: time.Now().Unix(),
		Body: []byte("hello world"),
	}

	p := randomPeerId()

	n.onMessage(g, p, m)

	sm, _ := n.ms.GetMessage(m.ID())
	if !reflect.DeepEqual(sm, m) {
		t.Errorf("message was not stored correctly")
	}

	s := n.state(g, m.ID(), p)

	if s.HoldFlag != true || s.AckFlag != true {
		t.Errorf("did not set flags")
	}
}

func getNodeForMessageHandlerTest() Node {
	n := Node{}
	n.syncState = make(map[GroupID]map[MessageID]map[PeerId]*State)
	n.offeredMessages = make(map[GroupID]map[PeerId][]MessageID)
	n.ID = randomPeerId()
	return n
}

func randomMessageId() MessageID {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	id := MessageID{}
	copy(id[:], bytes)

	return id
}

func randomPeerId() PeerId {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return PeerId(key.PublicKey)
}