package mvds

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"log"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	m.Run()
}

func TestNode_IsPeerInGroup_True(t *testing.T) {
	n := Node{}

	g := groupId("foo")

	peers := []PeerId{peerId(), peerId(), peerId()}

	n.sharing = map[GroupID][]PeerId{
		g: peers,
	}

	for i, tt := range peers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if !n.IsPeerInGroup(g, tt) {
				t.Error("expected peer is not in group")
			}
		})
	}
}

func TestNode_IsPeerInGroup_False(t *testing.T) {
	n := Node{}

	g := groupId("foo")

	peers := []PeerId{peerId(), peerId(), peerId()}

	n.sharing = map[GroupID][]PeerId{
		g: {},
	}

	for i, tt := range peers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if n.IsPeerInGroup(g, tt) {
				t.Error("expected peer is in group")
			}
		})
	}
}

func TestNode_onAck(t *testing.T) {
	ctrl := gomock.NewController(t)

	state := NewMocksyncState(ctrl)

	n := Node{}
	n.ID = peerId()
	n.syncState = state

	group := groupId("foo")
	peer := peerId()
	id := []byte("test")

	state.
		EXPECT().
		Remove(group, toMessageID(id), peer)

	ack := Ack{Id: [][]byte{id}}

	n.onAck(group, peer, ack)
}

func groupId(n string) GroupID {
	bytes := []byte(n)
	id := GroupID{}
	copy(id[:], bytes)
	return id
}

func peerId() PeerId {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return PeerId(key.PublicKey)
}
