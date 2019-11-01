package dependency

import (
	"sync"

	"github.com/vacp2p/mvds/state"
)

// Verify that MessageDependency interface is implemented.
var _ MessageDependency = (*memoryDependency)(nil)

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

func (md *memoryDependency) Add(msg, dependency state.MessageID) error {
	md.Lock()
	defer md.Unlock()
	// @todo check it wasn't already added
	md.dependents[dependency] = append(md.dependents[dependency], msg)
	md.dependencies[msg] += 1
	return nil
}

func (md *memoryDependency) Dependants(id state.MessageID) ([]state.MessageID, error) {
	md.Lock()
	defer md.Unlock()

	return md.dependents[id], nil
}

func (md *memoryDependency) MarkResolved(msg state.MessageID, dependency state.MessageID) error {
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
		return nil
	}

	md.dependents[dependency] = remove(md.dependents[dependency], id)
	md.dependencies[msg] -= 1
	return nil
}

func (md *memoryDependency) HasUnresolvedDependencies(id state.MessageID) (bool, error) {
	md.Lock()
	defer md.Unlock()

	return len(md.dependencies) > 0, nil
}

func remove(s []state.MessageID, i int) []state.MessageID {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
