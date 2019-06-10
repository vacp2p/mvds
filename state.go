package mvds

type State struct {
	SendCount   uint64
	SendEpoch   int64
}

type SyncState interface {
	Get(group GroupID, id MessageID, sender PeerID) State
	Set(group GroupID, id MessageID, sender PeerID, newState State)
	Remove(group GroupID, id MessageID, sender PeerID)
	Map(process func(GroupID, MessageID, PeerID, State) State)
}
