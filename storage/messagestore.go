package storage

import "github.com/status-im/mvds"

type MessageStore interface {
	HasMessage(id mvds.MessageID) bool
	GetMessage(id mvds.MessageID) (mvds.Message, error)
	SaveMessage(message mvds.Message) error
}
