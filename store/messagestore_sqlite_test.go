package store

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/vacp2p/mvds/state"

	"github.com/stretchr/testify/require"
	"github.com/vacp2p/mvds/persistenceutil"
	"github.com/vacp2p/mvds/protobuf"
	"github.com/vacp2p/mvds/store/migrations"
)

func TestPersistentMessageStore(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	db, err := persistenceutil.Open(tmpFile.Name(), "", persistenceutil.MigrationConfig{
		AssetNames:  migrations.AssetNames(),
		AssetGetter: migrations.Asset,
	})
	require.NoError(t, err)
	p := NewPersistentMessageStore(db)

	now := time.Now().Unix()
	message := protobuf.Message{
		GroupId:   []byte{0x01},
		Timestamp: now,
		Body:      []byte{0xaa, 0xbb, 0xcc},
	}

	err = p.Add(&message)
	require.NoError(t, err)
	// Adding the same message twice is not allowed.
	err = p.Add(&message)
	require.EqualError(t, err, "UNIQUE constraint failed: mvds_messages.id")
	// Verify if saved.
	exists, err := p.Has(message.ID())
	require.NoError(t, err)
	require.True(t, exists)
	recvMessage, err := p.Get(message.ID())
	require.NoError(t, err)
	require.Equal(t, message, *recvMessage)

	// Verify methods against non existing message.
	recvMessage, err = p.Get(state.MessageID{0xff})
	require.EqualError(t, err, "sql: no rows in result set")
	require.Nil(t, recvMessage)
	exists, err = p.Has(state.MessageID{0xff})
	require.NoError(t, err)
	require.False(t, exists)
}
