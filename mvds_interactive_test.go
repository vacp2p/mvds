package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vacp2p/mvds/node"
	"github.com/vacp2p/mvds/peers"
	"github.com/vacp2p/mvds/state"
	"github.com/vacp2p/mvds/store"
	"github.com/vacp2p/mvds/transport"
	"go.uber.org/zap"
)

func TestMVDSInteractiveSuite(t *testing.T) {
	suite.Run(t, new(MVDSInteractiveSuite))
}

type MVDSInteractiveSuite struct {
	suite.Suite
	client1 *node.Node
	client2 *node.Node
	ds1     store.MessageStore
	ds2     store.MessageStore
	state1  state.SyncState
	state2  state.SyncState
	peers1  peers.Persistence
	peers2  peers.Persistence
	groupID state.GroupID
}

func (s *MVDSInteractiveSuite) SetupTest() {

	logger := zap.NewNop()

	in1 := make(chan transport.Packet)
	t1 := transport.NewChannelTransport(0, in1)
	s.ds1 = store.NewDummyStore()
	s.state1 = state.NewSyncState()
	s.peers1 = peers.NewMemoryPersistence()
	p1 := [65]byte{0x01}
	s.client1 = node.NewNode(s.ds1, t1, s.state1, Calc, 0, p1, node.INTERACTIVE, s.peers1, logger)

	in2 := make(chan transport.Packet)
	t2 := transport.NewChannelTransport(0, in2)
	s.ds2 = store.NewDummyStore()
	s.state2 = state.NewSyncState()
	p2 := [65]byte{0x02}
	s.peers2 = peers.NewMemoryPersistence()
	s.client2 = node.NewNode(s.ds2, t2, s.state2, Calc, 0, p2, node.INTERACTIVE, s.peers2, logger)

	t2.AddOutput(p1, in1)
	t1.AddOutput(p2, in2)

	s.groupID = [32]byte{0x01, 0x2, 0x3, 0x4}

	s.Require().NoError(s.client1.AddPeer(s.groupID, p2))
	s.Require().NoError(s.client2.AddPeer(s.groupID, p1))

	s.client1.Start(10 * time.Millisecond)
	s.client2.Start(10 * time.Millisecond)
}

func (s *MVDSInteractiveSuite) TearDownTest() {
	s.client1.Stop()
	s.client2.Stop()
}

func (s *MVDSInteractiveSuite) TestInteractiveMode() {
	subscription := s.client2.Subscribe()
	messageID, err := s.client1.AppendMessage(s.groupID, []byte("message 1"))
	s.Require().NoError(err)

	// Check message is in store
	message1Sender, err := s.ds1.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Sender)

	// Check state is updated correctly
	states, err := s.state1.All(s.client1.CurrentEpoch())
	s.Require().NoError(err)
	s.Require().Equal(1, len(states))

	// Check we store the request
	s.Require().Eventually(func() bool {
		states, err := s.state2.All(s.client2.CurrentEpoch())
		return err == nil && len(states) == 1 && states[0].Type == state.REQUEST
	}, 1*time.Second, 10*time.Millisecond, "An request is stored in the state")

	<-subscription
	message1Receiver, err := s.ds2.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Receiver)

	// Check state is removed
	s.Require().Eventually(func() bool {
		states, err := s.state1.All(s.client1.CurrentEpoch())
		return err == nil && len(states) == 0

	}, 1*time.Second, 10*time.Millisecond, "We clear all the state")
}
