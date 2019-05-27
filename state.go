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

	// @todo: `Get` shouldn't mutate state, it should just check it
	//        this code should be in `Set`, here just check the state.
	//        that way you will also be able to use s.RLock() and s.RUnlock() here,
	//        you have an RWMutex anyway.
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

// @todo replace with somethind like
// Map(mutateFunc func(GroupID, MessageID, PeerID, state) state)
// so you iterate over the map there and set the value if needed
//
// your code will roughly look like
//
// s.Lock()
// defer s.Unlock()
// for groupID, groupMessages := range self.state {
//    for messageID, peerMessages := range groupMessages {
//       for peerID, state := range peerMessages {
//           s.state[groupID][messageID][peerID] = mutateFunc(groupID, messageID, peerID, state)
//       }
//    }
// }
//
func (s syncState) Iterate() map[GroupID]map[MessageID]map[PeerId]state {
	return s.state
}
