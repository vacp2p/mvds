package store

import (
	"errors"
	"sync"

	"github.com/vacp2p/mvds/protobuf"
	"github.com/vacp2p/mvds/state"
)

type DummyStore struct {
	sync.Mutex
	ms map[state.MessageID]*protobuf.Message
}

func NewDummyStore() *DummyStore {
	return &DummyStore{ms: make(map[state.MessageID]*protobuf.Message)}
}

func (ds *DummyStore) Has(id state.MessageID) (bool, error) {
	ds.Lock()
	defer ds.Unlock()

	_, ok := ds.ms[id]
	return ok, nil
}

func (ds *DummyStore) Get(id state.MessageID) (*protobuf.Message, error) {
	ds.Lock()
	defer ds.Unlock()

	m, ok := ds.ms[id]
	if !ok {
		return nil, errors.New("message does not exist")
	}

	return m, nil
}

func (ds *DummyStore) Add(message *protobuf.Message) error {
	ds.Lock()
	defer ds.Unlock()
	ds.ms[message.ID()] = message
	return nil
}

func (ds *DummyStore) GetMessagesWithoutChildren(group state.GroupID) ([]state.MessageID, error) {
	hasChildren := make(map[state.MessageID]bool)

	for id, msg := range ds.ms {
		if state.ToGroupID(msg.GroupId) != group {
			continue
		}

		if msg.Metadata != nil {
			for _, parent := range msg.Metadata.Parents {
				hasChildren[state.ToMessageID(parent)] = true
			}
		}

		if hasChildren[id] {
			continue
		}

		hasChildren[id] = false
	}

	msgs := make([]state.MessageID, 0)
	for id, hasChildren := range hasChildren {
		if hasChildren {
			continue
		}

		msgs = append(msgs, id)
	}

	return msgs, nil
}
