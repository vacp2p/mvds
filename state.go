package mvds

import "sync"

type state struct {
	HoldFlag    bool
	AckFlag     bool
	RequestFlag bool
	SendCount   uint64
	SendEpoch   int64
}

type syncState struct {
	sync.RWMutex

	state map[GroupID]map[MessageID]map[PeerId]state
}

func (s syncState) Get(group GroupID, id MessageID, sender PeerId) state {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.state[group]; !ok {
		s.state[group] = make(map[MessageID]map[PeerId]state)
	}

	if _, ok := s.state[group][id]; !ok {
		s.state[group][id] = make(map[PeerId]state)
	}

	if _, ok := s.state[group][id][sender]; !ok {
		s.state[group][id][sender] = state{}
	}

	return s.state[group][id][sender]
}

func (s *syncState) Set(group GroupID, id MessageID, sender PeerId, state state) {
	s.Lock()
	defer s.Unlock()

	s.state[group][id][sender] = state
}

func (s syncState) Iterate() map[GroupID]map[MessageID]map[PeerId]state {
	return s.state
}

func (s syncState) Map(callback func(g GroupID, m MessageID, p PeerId, s state) *state) {
	for group, syncstate := range s.state {
		for id, peers := range syncstate {
			for peer, state := range peers {
				newState := callback(group, id, peer, state)
				if newState == nil {
					continue
				}

				s.state[group][id][peer] = *newState
			}
		}
	}
}
