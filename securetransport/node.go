package securetransport

import "github.com/status-im/mvds"

type Node interface {
	SendPayload(senderId mvds.PeerId, to mvds.PeerId, payload mvds.Payload) error
}
