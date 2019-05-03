package mvds

type MessageStore interface {
	HasMessage(id MessageID) bool
	GetMessage(id MessageID) (Message, error)
	SaveMessage(message Message) error
}

type DummyStore struct {
	ms map[MessageID]Message
}

func (ds *DummyStore) HasMessage(id MessageID) bool {
	_, ok := ds.ms[id];  return ok
}

func (ds *DummyStore) GetMessage(id MessageID) (Message, error) {
	m, _ := ds.ms[id]; return m, nil
}

func (ds *DummyStore) SaveMessage(message Message) error {
	ds.ms[message.ID()] = message
	return nil
}
