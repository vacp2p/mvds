package mvds

import (
	"errors"
	"sync"
)

type DummyStore struct {
	sync.Mutex
	ms map[MessageID]Message
}

func NewDummyStore() DummyStore {
	return DummyStore{ms: make(map[MessageID]Message)}
}

func (ds *DummyStore) Has(id MessageID) bool {
	ds.Lock()
	defer ds.Unlock()

	_, ok := ds.ms[id]; return ok
}

func (ds *DummyStore) Get(id MessageID) (Message, error) {
	ds.Lock()
	defer ds.Unlock()

	m, ok := ds.ms[id]
	if !ok {
		return Message{}, errors.New("message does not exist")
	}

	return m, nil
}

func (ds *DummyStore) Add(message Message) error {
	ds.Lock()
	defer ds.Unlock()
	ds.ms[message.ID()] = message
	return nil
}
