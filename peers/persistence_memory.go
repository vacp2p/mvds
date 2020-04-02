package peers

import "github.com/vacp2p/mvds/state"

type memoryPersistence struct {
	peers map[state.GroupID][]state.PeerID
}

func NewMemoryPersistence() *memoryPersistence {
	return &memoryPersistence{
		peers: make(map[state.GroupID][]state.PeerID),
	}
}

func (p *memoryPersistence) Add(groupID state.GroupID, peerID state.PeerID) error {
	p.peers[groupID] = append(p.peers[groupID], peerID)
	return nil
}

func (p *memoryPersistence) Exists(groupID state.GroupID, peerID state.PeerID) (bool, error) {
	for _, peer := range p.peers[groupID] {
		if peer == peerID {
			return true, nil
		}
	}
	return false, nil
}

func (p *memoryPersistence) GetByGroupID(groupID state.GroupID) ([]state.PeerID, error) {
	return p.peers[groupID], nil
}

