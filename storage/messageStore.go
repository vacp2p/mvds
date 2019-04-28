package storage

import "github.com/status-im/mvds"

// @todo figure out how we add sync states
type SynchronisationState string

const (
	HOLD_FLAG    SynchronisationState = "hold"
	ACK_FLAG                          = "ack"
	REQUEST_FLAG                      = "request"
)

type MessageStore interface {
	GetMessage(id mvds.MessageID) (mvds.Message, error)
	SaveMessage(message mvds.Message) error
}
