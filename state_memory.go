package mvds

import "sync"

type memorySyncState struct {
	sync.Mutex

	state map[GroupID]map[MessageID]map[PeerId]state
}

func newSyncState() *memorySyncState {
	return &memorySyncState{
		state: make(map[GroupID]map[MessageID]map[PeerId]state),
	}
}

func (s *memorySyncState) Get(group GroupID, id MessageID, sender PeerId) state {
	s.Lock()
	defer s.Unlock()

	state, _ := s.state[group][id][sender]
	return state
}

func (s *memorySyncState) Set(group GroupID, id MessageID, sender PeerId, newState state) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.state[group]; !ok {
		s.state[group] = make(map[MessageID]map[PeerId]state)
	}

	if _, ok := s.state[group][id]; !ok {
		s.state[group][id] = make(map[PeerId]state)
	}

	s.state[group][id][sender] = newState
}

func (s *memorySyncState) Remove(group GroupID, id MessageID, sender PeerId) {
	s.Lock()
	defer s.Unlock()

	delete(s.state[group][id], sender)
}

func (s *memorySyncState) Map(process func(g GroupID, m MessageID, p PeerId, s state) state) {
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