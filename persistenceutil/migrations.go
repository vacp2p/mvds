package persistenceutil

import (
	nodemigrations "github.com/vacp2p/mvds/node/migrations"
	peersmigrations "github.com/vacp2p/mvds/peers/migrations"
	statemigrations "github.com/vacp2p/mvds/state/migrations"
	storemigrations "github.com/vacp2p/mvds/store/migrations"
)

type Migration struct {
	Names  []string
	Getter func(name string) ([]byte, error)
}

// DefaultMigrations is a collection of all mvds components migrations.
var DefaultMigrations = map[string]Migration{
	"node": {
		Names:  nodemigrations.AssetNames(),
		Getter: nodemigrations.Asset,
	},
	"peers": {
		Names:  peersmigrations.AssetNames(),
		Getter: peersmigrations.Asset,
	},
	"state": {
		Names:  statemigrations.AssetNames(),
		Getter: statemigrations.Asset,
	},
	"store": {
		Names:  storemigrations.AssetNames(),
		Getter: storemigrations.Asset,
	},
}
