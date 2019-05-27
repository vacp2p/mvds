package mvds

// @todo: rename this file to `sync_messageid.go` to show that it is related to
// sync.pb.go
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
