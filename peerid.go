package mvds

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/crypto"
)

type PeerId [64]byte

func PublicKeyToPeerID(k ecdsa.PublicKey) *PeerId {
	if k.X == nil || k.Y == nil {
		return nil
	}

	p := &PeerId{}
	copy(p[:], elliptic.Marshal(crypto.S256(), k.X, k.Y))
	return p
}
