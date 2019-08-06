package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vacp2p/mvds/node"
	"github.com/vacp2p/mvds/peers"
	"github.com/vacp2p/mvds/protobuf"
	"github.com/vacp2p/mvds/state"
	"github.com/vacp2p/mvds/store"
	"github.com/vacp2p/mvds/transport"
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

	in1 := make(chan transport.Packet)
	t1 := transport.NewChannelTransport(0, in1)
	ds1 := store.NewDummyStore()
	s.ds1 = &ds1
	s.state1 = state.NewSyncState()
	s.peers1 = peers.NewMemoryPersistence()
	p1 := [65]byte{0x01}
	s.client1 = node.NewNode(s.ds1, t1, s.state1, Calc, 0, p1, node.BATCH, s.peers1)

	in2 := make(chan transport.Packet)
	t2 := transport.NewChannelTransport(0, in2)
	ds2 := store.NewDummyStore()
	s.ds2 = &ds2
	s.state2 = state.NewSyncState()
	p2 := [65]byte{0x02}
	s.peers2 = peers.NewMemoryPersistence()
	s.client2 = node.NewNode(s.ds2, t2, s.state2, Calc, 0, p2, node.BATCH, s.peers2)

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
	subscription := make(chan protobuf.Message)
	s.client2.Subscribe(subscription)
	content := []byte("message 1")

	messageID, err := s.client1.AppendMessage(s.groupID, content)
	s.Require().NoError(err)

	// Check message is in store
	message1Sender, err := s.ds1.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Sender)

	// Check message is received
	s.Require().Eventually(func() bool {
		message1Receiver, err := s.ds2.Get(messageID)
		return err == nil && message1Receiver != nil
	}, 1*time.Second, 10*time.Millisecond)

	message := <- subscription
	s.Equal(message.Body, content)
}

func (s *MVDSBatchSuite) TestSendClient2ToClient1() {
	subscription := make(chan protobuf.Message)
	s.client1.Subscribe(subscription)
	content := []byte("message 1")

	messageID, err := s.client2.AppendMessage(s.groupID, content)
	s.Require().NoError(err)

	// Check message is in store
	message1Sender, err := s.ds2.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Sender)

	// Check message is received
	s.Require().Eventually(func() bool {
		message1Receiver, err := s.ds1.Get(messageID)
		return err == nil && message1Receiver != nil
	}, 1*time.Second, 10*time.Millisecond)

	message := <- subscription
	s.Equal(message.Body, content)
}

func (s *MVDSBatchSuite) TestAcks() {
	messageID, err := s.client1.AppendMessage(s.groupID, []byte("message 1"))
	s.Require().NoError(err)

	// Check message is in store
	message1Sender, err := s.ds1.Get(messageID)
	s.Require().NoError(err)
	s.Require().NotNil(message1Sender)

	// Check state is updated correctly
	states, err := s.state1.All()
	s.Require().NoError(err)
	s.Require().Equal(1, len(states))

	// Check message is received
	s.Require().Eventually(func() bool {
		message1Receiver, err := s.ds2.Get(messageID)
		return err == nil && message1Receiver != nil
	}, 1*time.Second, 10*time.Millisecond)

	// Check state is removed
	s.Require().Eventually(func() bool {
		states, err := s.state1.All()
		return err == nil && len(states) == 0

	}, 1*time.Second, 10*time.Millisecond)
}
