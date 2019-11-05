package state

type MessageID [32]byte
type GroupID [32]byte

// ToMessageID converts a byte array to a MessageID.
func ToMessageID(b []byte) MessageID {
	var id MessageID
	copy(id[:], b)
	return id
}

// ToGroupID converts a byte array to a GroupID.
func ToGroupID(b []byte) GroupID {
	var id GroupID
	copy(id[:], b)
	return id
}
