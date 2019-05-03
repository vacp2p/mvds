package mvds

type MessageStore interface {
	HasMessage(id MessageID) bool
	GetMessage(id MessageID) (Message, error)
	SaveMessage(message Message) error
}

type dummystore struct {
	MS map[MessageID]Message
}

func (ds *dummystore) HasMessage(id MessageID) bool {
	_, ok := ds.MS[id];  return ok
}

func (ds *dummystore) GetMessage(id MessageID) (Message, error) {
	m, _ := ds.MS[id]; return m, nil
}

func (ds *dummystore) SaveMessage(message Message) error {
	ds.MS[message.ID()] = message
	return nil
}
