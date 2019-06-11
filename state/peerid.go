package state

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type PeerID [65]byte

// Turns an ECSDA PublicKey to a PeerID
func PublicKeyToPeerID(k ecdsa.PublicKey) PeerID {
	var p PeerID
	copy(p[:], crypto.FromECDSAPub(&k))
	return p
}
