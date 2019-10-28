package dependency

import (
	"github.com/vacp2p/mvds/state"
)

type memoryDependency struct {
	dependents map[state.MessageID][]state.MessageID
	dependencies map[state.MessageID]int

}

func (md *memoryDependency) Add(msg, dependency state.MessageID) {
	// @todo check it wasn't already added
	md.dependents[dependency] = append(md.dependents[dependency], msg)
	md.dependencies[msg] += 1
}

func (md *memoryDependency) Dependants(id state.MessageID) []state.MessageID {
	return md.dependents[id]
}

func (md *memoryDependency) MarkResolved(msg state.MessageID, dependency state.MessageID) {
	// @todo remove from array
	panic("implement me")
}

func (md *memoryDependency) HasUnresolvedDependencies(id state.MessageID) bool {
	return len(md.dependencies) > 0
}
