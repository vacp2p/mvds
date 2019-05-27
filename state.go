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
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.state[group][id][sender]; !ok {
		return state{}
	}

	return s.state[group][id][sender]
}

func (s *syncState) Set(group GroupID, id MessageID, sender PeerId, newState state) {
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

func (s syncState) Iterate() map[GroupID]map[MessageID]map[PeerId]state {
	return s.state
}

func (s syncState) Map(process func(g GroupID, m MessageID, p PeerId, s state) state) {
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
