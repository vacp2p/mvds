package mvds

import (
	"sync"
)

type memorySyncState struct {
	sync.Mutex

	state map[GroupID]map[MessageID]map[PeerID]state
}

func NewSyncState() *memorySyncState {
	return &memorySyncState{
		state: make(map[GroupID]map[MessageID]map[PeerID]state),
	}
}

func (s *memorySyncState) Get(group GroupID, id MessageID, sender PeerID) state {
	s.Lock()
	defer s.Unlock()

	state, _ := s.state[group][id][sender]
	return state
}

func (s *memorySyncState) Set(group GroupID, id MessageID, sender PeerID, newState state) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.state[group]; !ok {
		s.state[group] = make(map[MessageID]map[PeerID]state)
	}

	if _, ok := s.state[group][id]; !ok {
		s.state[group][id] = make(map[PeerID]state)
	}

	s.state[group][id][sender] = newState
}

func (s *memorySyncState) Remove(group GroupID, id MessageID, sender PeerID) {
	s.Lock()
	defer s.Unlock()

	delete(s.state[group][id], sender)
}

func (s *memorySyncState) Map(process func(g GroupID, m MessageID, p PeerID, s state) state) {
	s.Lock()
	defer s.Unlock()

	for group, syncstate := range s.state {
		for id, peers := range syncstate {
			for peer, state := range peers {
				s.state[group][id][peer] = process(group, id, peer, state)
			}
		}
	}
}