package mvds

type MessageStore interface {
	HasMessage(id MessageID) bool
	GetMessage(id MessageID) (Message, error)
	SaveMessage(message Message) error
}
