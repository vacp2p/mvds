package mvds

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/crypto"
)

type PeerID [64]byte

func PublicKeyToPeerID(k ecdsa.PublicKey) PeerID {
	var p PeerID

	if k.X == nil || k.Y == nil {
		return p
	}

	copy(p[:], elliptic.Marshal(crypto.S256(), k.X, k.Y))
	return p
}
