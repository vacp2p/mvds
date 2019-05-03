package mvds

import (
	"crypto/rand"
	"reflect"
	"testing"
	"time"
)

func TestOnRequest(t *testing.T) {
	m := randomMessageId()

	r := Request{}
	r.Id = append(r.Id, m[:])

	n := getNodeForMessageHandlerTest()

	p := randomPeerId()
	n.onRequest(p, r)

	if n.state(m, p).RequestFlag != true {
		t.Errorf("did not set Request flag to true")
	}
}

func TestOnAck(t *testing.T) {
	m := randomMessageId()

	a := Ack{}
	a.Id = append(a.Id, m[:])

	n := getNodeForMessageHandlerTest()

	p := randomPeerId()
	n.onAck(p, a)

	if n.state(m, p).HoldFlag != true {
		t.Errorf("did not set Hold flag to true")
	}
}

func TestOnOffer(t *testing.T) {
	m := randomMessageId()

	o := Offer{}
	o.Id = append(o.Id, m[:])

	n := getNodeForMessageHandlerTest()

	p := randomPeerId()
	n.onOffer(p, o)

	if n.state(m, p).HoldFlag != true {
		t.Errorf("did not set Hold flag to true")
	}

	if n.offeredMessages[p][0] != m {
		t.Errorf("message was not added to offered list")
	}
}

func TestOnMessage(t *testing.T) {
	n := getNodeForMessageHandlerTest()

	ds := NewDummyStore()
	n.ms = &ds

	id := randomMessageId()

	m := Message{
		GroupId: id[:],
		Timestamp: time.Now().Unix(),
		Body: []byte("hello world"),
	}

	p := randomPeerId()

	n.onMessage(p, m)

	sm, _ := n.ms.GetMessage(m.ID())
	if !reflect.DeepEqual(sm, m) {
		t.Errorf("message was not stored correctly")
	}

	s := n.state(m.ID(), p)

	if s.HoldFlag != true || s.AckFlag != true {
		t.Errorf("did not set flags")
	}
}

func getNodeForMessageHandlerTest() Node {
	n := Node{}
	n.syncState = make(map[MessageID]map[PeerId]*State)
	n.offeredMessages = make(map[PeerId][]MessageID)
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
	bytes := make([]byte, 32)
	rand.Read(bytes)

	id := PeerId{}
	copy(id[:], bytes)

	return id
}