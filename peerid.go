package mvds

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type PeerID [64]byte

func PublicKeyToPeerID(k ecdsa.PublicKey) PeerID {
	var p PeerID
	copy(p[:], crypto.FromECDSAPub(&k))
	return p
}
