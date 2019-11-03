package dependency

import (
	"sync"

	"github.com/vacp2p/mvds/state"
)

// Verify that Tracker interface is implemented.
var _ Tracker = (*inMemoryTracker)(nil)

type inMemoryTracker struct {
	sync.Mutex

	dependents map[state.MessageID][]state.MessageID
	dependencies map[state.MessageID]int

}

func NewInMemoryTracker() *inMemoryTracker {
	return &inMemoryTracker{
		dependents: make(map[state.MessageID][]state.MessageID),
		dependencies: make(map[state.MessageID]int),
	}
}

func (md *inMemoryTracker) Add(msg, dependency state.MessageID) error {
	md.Lock()
	defer md.Unlock()
	// @todo check it wasn't already added
	md.dependents[dependency] = append(md.dependents[dependency], msg)
	md.dependencies[msg] += 1
	return nil
}

func (md *inMemoryTracker) Dependants(id state.MessageID) ([]state.MessageID, error) {
	md.Lock()
	defer md.Unlock()

	return md.dependents[id], nil
}

func (md *inMemoryTracker) MarkResolved(msg state.MessageID, dependency state.MessageID) error {
	md.Lock()
	defer md.Unlock()

	for i, item := range md.dependents[dependency] {
		if item != msg {
			continue
		}

		md.dependents[dependency] = remove(md.dependents[dependency], i)
		md.dependencies[msg] -= 1
		return nil
	}

	return nil
}

func (md *inMemoryTracker) HasUnresolvedDependencies(id state.MessageID) (bool, error) {
	md.Lock()
	defer md.Unlock()

	return len(md.dependencies) > 0, nil
}

func remove(s []state.MessageID, i int) []state.MessageID {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
