package mvds

type state struct {
	SendCount   uint64
	SendEpoch   uint64
}

type syncState interface {
	Get(group GroupID, id MessageID, sender PeerID) state
	Set(group GroupID, id MessageID, sender PeerID, newState state)
	Remove(group GroupID, id MessageID, sender PeerID)
	Map(process func(g GroupID, m MessageID, p PeerID, s state) state)
}
