package mvds

import "github.com/status-im/mvds/protobuf"

type MessageStore interface {
	Has(id MessageID) bool
	Get(id MessageID) (protobuf.Message, error)
	Add(message protobuf.Message) error
}
