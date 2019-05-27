package mvds

import (
	"sync"
)

type MessageStore interface {
	// @todo: `Add`, `Get` and `Has` is a good pattern to follow
	HasMessage(id MessageID) bool
	GetMessage(id MessageID) (Message, error)
	SaveMessage(message Message) error
}

// @todo I would extact it to `messagestore_dummy.go`
type DummyStore struct {
	sync.Mutex
	ms map[MessageID]Message
}

func NewDummyStore() DummyStore {
	return DummyStore{ms: make(map[MessageID]Message)}
}

func (ds *DummyStore) HasMessage(id MessageID) bool {
	ds.Lock()
	defer ds.Unlock()

	_, ok := ds.ms[id]
	return ok
}

func (ds *DummyStore) GetMessage(id MessageID) (Message, error) {
	ds.Lock()
	defer ds.Unlock()

	m, _ := ds.ms[id]
	return m, nil
}

func (ds *DummyStore) SaveMessage(message Message) error {
	ds.Lock()
	defer ds.Unlock()
	ds.ms[message.ID()] = message
	return nil
}
