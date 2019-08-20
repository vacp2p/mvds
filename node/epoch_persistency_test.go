package node

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vacp2p/mvds/node/migrations"
	"github.com/vacp2p/mvds/persistenceutil"
	"github.com/vacp2p/mvds/state"
)

func TestEpochSQLitePersistence(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	db, err := persistenceutil.Open(tmpFile.Name(), "", persistenceutil.MigrationConfig{
		AssetNames:  migrations.AssetNames(),
		AssetGetter: migrations.Asset,
	})
	require.NoError(t, err)
	p := newEpochSQLitePersistence(db)

	err = p.Set(state.PeerID{0x01}, 1)
	require.NoError(t, err)
	epoch, err := p.Get(state.PeerID{0x01})
	require.NoError(t, err)
	require.Equal(t, int64(1), epoch)

	err = p.Set(state.PeerID{0x01}, 2)
	require.NoError(t, err)
	epoch, err = p.Get(state.PeerID{0x01})
	require.NoError(t, err)
	require.Equal(t, int64(2), epoch)

	epoch, err = p.Get(state.PeerID{0xff})
	require.NoError(t, err)
	require.Equal(t, int64(0), epoch)
}
