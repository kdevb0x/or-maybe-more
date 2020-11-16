package client

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/stretchr/objx"
	"golang.org/x/crypto/nacl/auth"
)

type clientID struct {
	dna         ed25519.PrivateKey
	displayName string
}

// An clientID is constructed by initializing its DNA.
func initDNA() *clientID {
	_, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(fmt.Errorf("failed to generate keys: %w\n", err).Error())
	}
	return &clientID{dna: privKey}
}

// Verify that the message msg has been signed by this clientID.
func (id *clientID) Verify(msg []byte) (bool, error) {
	if id.dna == nil {
		return false, errors.New(fmt.Sprintf("DNA not initialized for %s\n", id.displayName))
	}
	ok := auth.Verify(msg, nil, privKeyToBP(id.dna))
	if !ok {
		return false, errors.New("Verification failed")
	}
	return true, nil
}

// Sign the message msg using this clientID.
func (id *clientID) Sign(msg []byte) (sig []byte, err error) {
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

// Client represents a user.
type Client struct {
	ipaddr   string
	metadata Metadata
}

// Metadata is used collected from clients, and used to help derive possible
// character traits of a user in order to gain a clearer picture of the
// individual for compatibility matching.
type Metadata interface {
	// Store iterates over m, saving the values to storage.
	Store(m MetadataQuery) error
	// Fetch returns the data requested in m (if present).
	Fetch(m MetadataQuery) (data /*FIXME*/ interface{}, err error)
}

// MetadataQuery is used to query, and fetch Client Metadata.
type MetadataQuery map[string]string

func (mq *MetadataQuery) Store() error {

}

// BaseMetadataQuery is the most basic MetadataQuery, which contains only "last_seen",
// and "display_name".
var BaseMetadataQuery = MetadataQuery{"last_seen": "", "display_name": ""}

type authdata struct {
	client *Client

	// auth data from provider
	authobj   objx.Map
	hasCookie bool
}
