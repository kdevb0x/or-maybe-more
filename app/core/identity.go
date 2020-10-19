package core

import (
	"crypto/ed25519"
)

type identity struct {
	dna ed25519.PrivateKey
}

func initDNA(seed []byte) *identity {
	privKey := ed25519.NewKeyFromSeed(seed)
	return &identity{dna: privKey}
}
