package dependency

import (
	"github.com/vacp2p/mvds/state"
)

type Tracker interface {
	Add(msg, dependency state.MessageID) error
	Dependants(id state.MessageID) ([]state.MessageID, error)
	Resolve(msg state.MessageID, dependency state.MessageID) error
	IsResolved(id state.MessageID) (bool, error)
}
