package core

import (
	"crypto/ed25519"

	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/auth"
)

type identity struct {
	dna         ed25519.PrivateKey
	displayName string
}

// An identity is constructed by initializing its DNA.
func initDNA() *identity {
	_, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Errorf("failed to generate keys: %w\n", err).Error())
	}
	return &identity{dna: privKey}
}

// Verify that the message msg has been signed by this identity.
func (id *identity) Verify(msg []byte) (bool, error) {
	if id.dna == nil {
		return false, errors.New(fmt.Sprintf("DNA not initialized for %s\n", id.displayName))
	}
	ok := auth.Verify(msg, nil, privKeyToBP(id.dna))
	if !ok {
		return false, errors.New("Verification failed")
	}
	return true, nil
}

// Sign the message msg using this identity.
func (id *identity) Sign(msg []byte) (sig []byte, err error) {
	if len(id.dna) == 0 {
		return nil, errors.New("error: the dna of this identidy has not been initialized")
	}
	return ed25519.Sign(id.dna, msg), nil
}

// Converts a byte slice to a byte array pointer.
func bToBP(b []byte) *[32]byte {
	var s *[32]byte
	// since k uses s as the underlying array, we this should work.
	k := s[:]
	copy(k, b[:31])
	return s
}

// Converts from ed25519.PrivateKey to a byte array pointer.
func privKeyToBP(k ed25519.PrivateKey) *[32]byte {
	b := []byte(k)
	return bToBP(b)
}
