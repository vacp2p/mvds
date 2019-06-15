// Package state contains everything related to the synchronization state for MVDS.
package state

import "github.com/status-im/mvds/epoch"

type State struct {
	SendCount   uint64
	SendEpoch   uint64
}

type SyncState interface {
	Get(group GroupID, id MessageID, peer PeerID) (State, error)
	Set(group GroupID, id MessageID, peer PeerID, newState State) error
	Remove(group GroupID, id MessageID, peer PeerID) error
	Map(epoch *epoch.Epoch, process func(GroupID, MessageID, PeerID, State) State) error
}
