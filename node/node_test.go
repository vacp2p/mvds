package node

import (
	"crypto/rand"
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

	peer := peerID()

	group := groupId()

	parent := &protobuf.Message{
		GroupId: group[:],
		Timestamp: time.Now().Unix(),
		Body: []byte{0x01},
	}

	parentID := parent.ID()

	msg := &protobuf.Message{
		GroupId: group[:],
		Timestamp: time.Now().Unix(),
		Body: []byte{0x01},
		Metadata: &protobuf.Metadata{Ephemeral: false, Parents: [][]byte{parentID[:]}},
	}

	syncstate.EXPECT().Add(gomock.Any()).Return(nil)

	node.resolveEventually(peer, msg)
}

func peerID() (id state.PeerID) {
	_, _ = rand.Read(id[:])
	return id
}

func groupId() (id state.GroupID) {
	_, _ = rand.Read(id[:])
	return id
}
