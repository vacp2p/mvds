package mvds

type MessageType int

const (
	OFFER   MessageType = 0
	REQUEST             = 1
)

type state struct {
	Type      MessageType
	SendCount uint64
	SendEpoch int64
}

type syncState interface {
	Get(group GroupID, id MessageID, sender PeerId) state
	Set(group GroupID, id MessageID, sender PeerId, newState state)
	Remove(group GroupID, id MessageID, sender PeerId)
	Map(process func(g GroupID, m MessageID, p PeerId, s state) state)
}
