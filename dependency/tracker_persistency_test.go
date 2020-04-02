package dependency

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vacp2p/mvds/dependency/migrations"
	"github.com/vacp2p/mvds/persistenceutil"
	"github.com/vacp2p/mvds/state"
)

func TestTrackerSQLitePersistence(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	db, err := persistenceutil.Open(tmpFile.Name(), "", persistenceutil.MigrationConfig{
		AssetNames:  migrations.AssetNames(),
		AssetGetter: migrations.Asset,
	})
	require.NoError(t, err)
	d := NewPersistentTracker(db)

	msg := state.MessageID{0x01}
	dependency := state.MessageID{0x02}

	err = d.Add(msg, dependency)
	require.NoError(t, err)
	dependants, err := d.Dependants(dependency)
	require.NoError(t, err)
	require.Equal(t, msg, dependants[0])

	res, err := d.IsResolved(msg)
	require.NoError(t, err)
	require.False(t, res)

	err = d.Resolve(msg, dependency)
	require.NoError(t, err)

	res, err = d.IsResolved(msg)
	require.NoError(t, err)
	require.True(t, res)
}
