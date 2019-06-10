package mvds

import (
	"errors"
	"sync"

	"github.com/status-im/mvds/protobuf"
)

type DummyStore struct {
	sync.Mutex
	ms map[MessageID]protobuf.Message
}

func NewDummyStore() DummyStore {
	return DummyStore{ms: make(map[MessageID]protobuf.Message)}
}

func (ds *DummyStore) Has(id MessageID) bool {
	ds.Lock()
	defer ds.Unlock()

	_, ok := ds.ms[id]; return ok
}

func (ds *DummyStore) Get(id MessageID) (protobuf.Message, error) {
	ds.Lock()
	defer ds.Unlock()

	m, ok := ds.ms[id]
	if !ok {
		return protobuf.Message{}, errors.New("message does not exist")
	}

	return m, nil
}

func (ds *DummyStore) Add(message protobuf.Message) error {
	ds.Lock()
	defer ds.Unlock()
	ds.ms[ID(message)] = message
	return nil
}
