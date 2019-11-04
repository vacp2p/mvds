package node

import (
	"crypto/rand"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/vacp2p/mvds/node/internal"
	"github.com/vacp2p/mvds/protobuf"
	"github.com/vacp2p/mvds/state"
	"github.com/vacp2p/mvds/store"
)


func TestNode_resolveEventually(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	syncstate := internal.NewMockSyncState(ctrl)

	node := Node{
		syncState: syncstate,
		store: store.NewMemoryMessageStore(),
	}

	channel := node.Subscribe()

	peer := peerID()
	group := groupID()
	parent := messageID()

	msg := &protobuf.Message{
		GroupId: group[:],
		Timestamp: time.Now().Unix(),
		Body: []byte{0x01},
		Metadata: &protobuf.Metadata{Ephemeral: false, Parents: [][]byte{parent[:]}},
	}

	expectedState := state.State{
		GroupID:   &group,
		MessageID: parent,
		PeerID:    peer,
		Type:      state.REQUEST,
		SendEpoch: 1,
	}

	syncstate.EXPECT().Add(expectedState).Return(nil)

	go node.resolveEventually(peer, msg)

	received := <-channel

	if !reflect.DeepEqual(*msg, received) {
		t.Error("expected message did not match received")
	}
}

func peerID() (id state.PeerID) {
	_, _ = rand.Read(id[:])
	return id
}

func groupID() (id state.GroupID) {
	_, _ = rand.Read(id[:])
	return id
}

func messageID() (id state.MessageID) {
	_, _ = rand.Read(id[:])
	return id
}
