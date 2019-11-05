package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vacp2p/mvds/dependency"
	"github.com/vacp2p/mvds/node"
	"github.com/vacp2p/mvds/peers"
	"github.com/vacp2p/mvds/state"
	"github.com/vacp2p/mvds/store"
	"github.com/vacp2p/mvds/transport"
	"go.uber.org/zap"
)

func TestMVDSBatchSuite(t *testing.T) {
	suite.Run(t, new(MVDSBatchSuite))
}

type MVDSBatchSuite struct {
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

func (s *MVDSBatchSuite) SetupTest() {

	logger := zap.NewNop()

	in1 := make(chan transport.Packet)
	t1 := transport.NewChannelTransport(0, in1)
	s.ds1 = store.NewMemoryMessageStore()
	s.state1 = state.NewMemorySyncState()
	s.peers1 = peers.NewMemoryPersistence()
	p1 := [65]byte{0x01}
	s.client1 = node.NewNode(s.ds1, t1, s.state1, Calc, 0, p1, node.BatchMode, s.peers1, dependency.NewInMemoryTracker(), node.EventualMode, logger)

	in2 := make(chan transport.Packet)
	t2 := transport.NewChannelTransport(0, in2)
	s.ds2 = store.NewMemoryMessageStore()
	s.state2 = state.NewMemorySyncState()
	p2 := [65]byte{0x02}
	s.peers2 = peers.NewMemoryPersistence()
	s.client2 = node.NewNode(s.ds2, t2, s.state2, Calc, 0, p2, node.BatchMode, s.peers2, dependency.NewInMemoryTracker(), node.EventualMode, logger)

	t2.AddOutput(p1, in1)
	t1.AddOutput(p2, in2)

	s.groupID = [32]byte{0x01, 0x2, 0x3, 0x4}

	s.Require().NoError(s.client1.AddPeer(s.groupID, p2))
	s.Require().NoError(s.client2.AddPeer(s.groupID, p1))

	// We run the tick manually
	s.client1.Start(10 * time.Millisecond)
	s.client2.Start(10 * time.Millisecond)
}

func (s *MVDSBatchSuite) TearDownTest() {
	s.client1.Stop()
	s.client2.Stop()
}

func (s *MVDSBatchSuite) TestSendClient1ToClient2() {
	subscription := s.client2.Subscribe()
	content := []byte("message 1")

	messageID, err := s.client1.AppendMessage(s.groupID, content)
	s.Require().NoError(err)

	// Check message is in store
	message1Sender, err := s.ds1.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Sender)

	message := <-subscription
	s.Equal(message.Body, content)

	message1Receiver, err := s.ds2.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Receiver)
}

func (s *MVDSBatchSuite) TestSendClient2ToClient1() {
	subscription := s.client1.Subscribe()
	content := []byte("message 1")

	messageID, err := s.client2.AppendMessage(s.groupID, content)
	s.Require().NoError(err)

	// Check message is in store
	message1Sender, err := s.ds2.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Sender)

	message := <-subscription
	s.Equal(message.Body, content)

	message1Receiver, err := s.ds1.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Receiver)
}

func (s *MVDSBatchSuite) TestAcks() {
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

	<-subscription

	message1Receiver, err := s.ds2.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Receiver)

	// Check state is removed
	s.Require().Eventually(func() bool {
		states, err := s.state1.All(s.client1.CurrentEpoch())
		return err == nil && len(states) == 0

	}, 1*time.Second, 10*time.Millisecond)
}
