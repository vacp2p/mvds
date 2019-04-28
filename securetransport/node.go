package securetransport

import "github.com/status-im/mvds"

type Node interface {
	SendMessage(senderId []byte, to []byte, message mvds.Message) error // @todo probably needs types change message to payload
}
