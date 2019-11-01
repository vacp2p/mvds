package dependency

import (
	"sync"

	"github.com/vacp2p/mvds/state"
)

type memoryDependency struct {
	sync.Mutex

	dependents map[state.MessageID][]state.MessageID
	dependencies map[state.MessageID]int

}

func NewDummyDependency() *memoryDependency {
	return &memoryDependency{
		dependents: make(map[state.MessageID][]state.MessageID),
		dependencies: make(map[state.MessageID]int),
	}
}

func (md *memoryDependency) Add(msg, dependency state.MessageID) {
	md.Lock()
	defer md.Unlock()
	// @todo check it wasn't already added
	md.dependents[dependency] = append(md.dependents[dependency], msg)
	md.dependencies[msg] += 1
}

func (md *memoryDependency) Dependants(id state.MessageID) []state.MessageID {
	md.Lock()
	defer md.Unlock()

	return md.dependents[id]
}

func (md *memoryDependency) MarkResolved(msg state.MessageID, dependency state.MessageID) {
	md.Lock()
	defer md.Unlock()

	id := -1
	for i, val := range md.dependents[dependency] {
		if val == msg {
			id = i
			break
		}
	}

	if id == -1 {
		return
	}

	md.dependents[dependency] = remove(md.dependents[dependency], id)
	md.dependencies[msg] -= 1
}

func (md *memoryDependency) HasUnresolvedDependencies(id state.MessageID) bool {
	md.Lock()
	defer md.Unlock()

	return len(md.dependencies) > 0
}

func remove(s []state.MessageID, i int) []state.MessageID {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
