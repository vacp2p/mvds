package state

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vacp2p/mvds/persistenceutil"
	"github.com/vacp2p/mvds/state/migrations"
)

func TestPersistentSyncState(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	db, err := persistenceutil.Open(tmpFile.Name(), "", persistenceutil.MigrationConfig{
		AssetNames:  migrations.AssetNames(),
		AssetGetter: migrations.Asset,
	})
	require.NoError(t, err)
	p := NewPersistentSyncState(db)

	stateWithoutGroupID := State{
		Type:      OFFER,
		SendCount: 1,
		SendEpoch: 123456,
		GroupID:   nil,
		PeerID:    PeerID{0x01},
		MessageID: MessageID{0xaa},
	}
	err = p.Add(stateWithoutGroupID)
	require.NoError(t, err)

	stateWithGroupID := stateWithoutGroupID
	stateWithGroupID.GroupID = &GroupID{0x01}
	stateWithGroupID.MessageID = MessageID{0xbb}
	err = p.Add(stateWithGroupID)
	require.NoError(t, err)

	allStates, err := p.All()
	require.NoError(t, err)
	require.Equal(t, []State{stateWithoutGroupID, stateWithGroupID}, allStates)
	require.Nil(t, allStates[0].GroupID)
	require.EqualValues(t, &GroupID{0x01}, allStates[1].GroupID)

	err = p.Remove(stateWithoutGroupID.MessageID, stateWithoutGroupID.PeerID)
	require.NoError(t, err)
	// remove non-existing row
	err = p.Remove(MessageID{0xff}, PeerID{0xff})
	require.EqualError(t, err, "state not found")
}
