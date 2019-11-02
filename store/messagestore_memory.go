package store

import (
	"errors"
	"sync"

	"github.com/vacp2p/mvds/protobuf"
	"github.com/vacp2p/mvds/state"
)

type memoryMessageStore struct {
	sync.Mutex
	ms map[state.MessageID]*protobuf.Message
}

func NewMemoryMessageStore() *memoryMessageStore {
	return &memoryMessageStore{ms: make(map[state.MessageID]*protobuf.Message)}
}

func (ds *memoryMessageStore) Has(id state.MessageID) (bool, error) {
	ds.Lock()
	defer ds.Unlock()

	_, ok := ds.ms[id]
	return ok, nil
}

func (ds *memoryMessageStore) Get(id state.MessageID) (*protobuf.Message, error) {
	ds.Lock()
	defer ds.Unlock()

	m, ok := ds.ms[id]
	if !ok {
		return nil, errors.New("message does not exist")
	}

	return m, nil
}

func (ds *memoryMessageStore) Add(message *protobuf.Message) error {
	ds.Lock()
	defer ds.Unlock()
	ds.ms[message.ID()] = message
	return nil
}
