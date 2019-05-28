package mvds

type MessageStore interface {
	Has(id MessageID) bool
	Get(id MessageID) (Message, error)
	Add(message Message) error
}
