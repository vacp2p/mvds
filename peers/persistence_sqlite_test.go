package peers

import (
	"io/ioutil"
	"testing"

	"github.com/vacp2p/mvds/peers/migrations"

	"github.com/vacp2p/mvds/persistenceutil"
	"github.com/vacp2p/mvds/state"

	"github.com/stretchr/testify/require"
)

func TestSQLitePersistence(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	db, err := persistenceutil.Open(tmpFile.Name(), "", persistenceutil.MigrationConfig{
		AssetNames:  migrations.AssetNames(),
		AssetGetter: migrations.Asset,
	})
	require.NoError(t, err)
	p := NewSQLitePersistence(db)

	err = p.Add(state.GroupID{0x01}, state.PeerID{0x01})
	require.NoError(t, err)
	// Add the same again.
	err = p.Add(state.GroupID{0x01}, state.PeerID{0x01})
	require.NoError(t, err)
	// Add another peer to the same group.
	err = p.Add(state.GroupID{0x01}, state.PeerID{0x02})
	require.NoError(t, err)
	// Create a new group.
	err = p.Add(state.GroupID{0x02}, state.PeerID{0x01})
	require.NoError(t, err)
	err = p.Add(state.GroupID{0x02}, state.PeerID{0x02})
	require.NoError(t, err)
	err = p.Add(state.GroupID{0x02}, state.PeerID{0x03})
	require.NoError(t, err)

	// Validate group 0x01.
	peers, err := p.GetByGroupID(state.GroupID{0x01})
	require.NoError(t, err)
	require.Equal(t, []state.PeerID{state.PeerID{0x01}, state.PeerID{0x02}}, peers)
	// Validate group 0x02.
	peers, err = p.GetByGroupID(state.GroupID{0x02})
	require.NoError(t, err)
	require.Equal(t, []state.PeerID{state.PeerID{0x01}, state.PeerID{0x02}, state.PeerID{0x03}}, peers)

	// Validate existence method.
	exists, err := p.Exists(state.GroupID{0x01}, state.PeerID{0x01})
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = p.Exists(state.GroupID{0x01}, state.PeerID{0xFF})
	require.NoError(t, err)
	require.False(t, exists)
}
