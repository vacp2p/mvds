package mvds

type state struct {
	SendCount   uint64
	SendEpoch   int64
}

type SyncState interface {
	Get(group GroupID, id MessageID, sender PeerID) state
	Set(group GroupID, id MessageID, sender PeerID, newState state)
	Remove(group GroupID, id MessageID, sender PeerID)
	Map(process func(g GroupID, m MessageID, p PeerID, s state) state)
}
