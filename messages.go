package mvds

import (
	"crypto/sha256"
	"encoding/binary"
)

type MessageID [32]byte
type GroupID [32]byte

func (m Message) ID() MessageID {
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(m.Timestamp))

	b := append([]byte("MESSAGE_ID"), m.GroupId[:]...)
	b = append(b, t...)
	b = append(b, m.Body...)

	return sha256.Sum256(b)
}
