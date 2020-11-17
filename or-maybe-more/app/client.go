package app

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/auth"
)

type ClientID struct {
	dna         ed25519.PrivateKey
	displayName string
	key         ed25519.PublicKey
	OAuthData   *OAuthUserData
}

// An ClientID is constructed by initializing its DNA.
func initClientDNA() *ClientID {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Errorf("failed to generate keys: %w\n", err).Error())
	}
	return &ClientID{dna: privKey, key: pubKey}
}

// Verify that the message msg has been signed by this ClientID.
func (id *ClientID) Verify(msg []byte) (ok bool, err error) {
	if (id.dna == nil || len(id.dna) == 0) || (id.key == nil || len(id.key) == 0) {
		return false, errors.New(fmt.Sprintf("DNA not initialized for %s\n", id.displayName))
	}
	ok = auth.Verify(id.key, msg, privKeyToBPtr(id.dna))
	if !ok {
		err = errors.New("Verification failed")
		return ok, err
	}
	return true, nil
}

// Sign the message msg using this ClientID.
func (id *ClientID) Sign(msg []byte) (sig []byte, err error) {
	if len(id.dna) == 0 {
		return nil, errors.New("error: the dna of this identidy has not been initialized")
	}
	return ed25519.Sign(id.dna, msg), nil
}

// Converts b[:31] to a 32 byte array pointer.
func bptr(b []byte) *[32]byte {
	var s *[32]byte
	// since k uses s as the underlying array, we this should work.
	k := s[:]
	copy(k, b[:31])
	return s
}

// Converts from ed25519.PrivateKey to a byte array pointer.
func privKeyToBPtr(k ed25519.PrivateKey) *[32]byte {
	b := []byte(k)
	return bptr(b)
}

// Client represents a user.
type Client struct {
	ipaddr   string
	Metadata Metadata
	id       *ClientID
}

// Metadata is collected from Clients, and used to help derive possible
// character traits of a user in order to gain a clearer picture of the
// individual for compatibility matching.
type Metadata interface {
	// Store iterates over m, saving the values to storage.
	Store(m MetadataQuery) error

	// Fetch ranges over m filling in the values of every key present,
	// returning the total number of changes made.
	Fetch(m MetadataQuery) (n int, err error)
}

// MetadataQuery is used to query, and fetch Client Metadata.
type MetadataQuery map[string]interface{}

func (mq *MetadataQuery) Store(m MetadataQuery) error {
	for k, v := range m {
		(*mq)[k] = v
	}
	return nil
}

func (mq *MetadataQuery) Fetch(m MetadataQuery) int {
	var cnt int
	for k := range m {
		if nv, ok := (*mq)[k]; ok {
			m[k] = nv
			cnt++
		}
	}
	return cnt
}

// BaseMetadataQuery is the most basic MetadataQuery, which contains only "last_seen",
// and "display_name".
var BaseMetadataQuery = MetadataQuery{"last_seen": "", "display_name": ""}
