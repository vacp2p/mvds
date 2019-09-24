package protobuf

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/multiformats/go-multihash"
	"github.com/vacp2p/mvds/state"
)

// ID creates the MessageID for a Message
func (m Message) ID() state.MessageID {
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(m.Timestamp))

	b := append([]byte("MESSAGE_ID"), m.GroupId[:]...)
	b = append(b, t...)
	b = append(b, m.Body...)

	hash := sha256.Sum256(b)
	id, _ := multihash.Encode(hash[:], multihash.SHA2_256)

	mid := state.MessageID{}
	copy(id[:], id)

	return mid
}
