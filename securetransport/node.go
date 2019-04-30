package securetransport

import "github.com/status-im/mvds"

type Node interface {
	SendPayload(senderId []byte, to []byte, payload mvds.Payload) error
}
