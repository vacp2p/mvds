package mvds

type state struct {
	SendCount   uint64
	SendEpoch   uint64
}

type syncState interface {
	Get(group GroupID, id MessageID, sender PeerId) state
	Set(group GroupID, id MessageID, sender PeerId, newState state)
	Remove(group GroupID, id MessageID, sender PeerId)
	Map(process func(g GroupID, m MessageID, p PeerId, s state) state)
}
