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

func (ds *DummyStore) GetMessagesWithoutChildren(id state.GroupID) []state.MessageID {
	msgs := make([]state.MessageID, 0)

	for msgid, msg := range ds.ms {
		if state.ToGroupID(msg.GroupId) != id {
			continue
		}

		if msg.Metadata == nil {
			continue
		}

		for _, parent := range msg.Metadata.Parents {
			for i, p := range msgs {
				if p == state.ToMessageID(parent) {
					msgs = remove(msgs, i)
					break
				}
			}
		}

		msgs = append(msgs, msgid)
	}

	return msgs
}

func remove(s []state.MessageID, i int) []state.MessageID {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
