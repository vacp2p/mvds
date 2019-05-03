package mvds

type MessageStore interface {
	HasMessage(id MessageID) bool
	GetMessage(id MessageID) (Message, error)
	SaveMessage(message Message) error
}

type dummystore struct {
	ms map[MessageID]Message
}

func (ds *dummystore) HasMessage(id MessageID) bool {
	_, ok := ds.ms[id];  return ok
}

func (ds *dummystore) GetMessage(id MessageID) (Message, error) {
	m, _ := ds.ms[id]; return m, nil
}

func (ds *dummystore) SaveMessage(message Message) error {
	ds.ms[message.ID()] = message
	return nil
}
