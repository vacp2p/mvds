package dependency

import (
	"github.com/vacp2p/mvds/state"
)

type MessageDependency interface {
	Add(msg, dependency state.MessageID)
	Dependants(id state.MessageID) []state.MessageID
	MarkResolved(msg state.MessageID, dependency state.MessageID)
	HasUnresolvedDependencies(id state.MessageID) bool
}
