package net

import (
	"github.com/stretchr/objx"
)

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
	Store(m MetaQuery) error
	// Fetch returns the data requested in m (if present).
	Fetch(m MetaQuery) (data /*FIXME*/ interface{}, err error)
}

// MetaQuery is used to query, and fetch Client Metadata.
type MetaQuery map[string]string

func (mq *MetaQuery) Store() error {

}

// BaseMetaQuery is the most basic MetaQuery, which contains only "last_seen",
// and "display_name".
var BaseMetaQuery = func() MetaQuery {
	m := make(map[string]string)
	m["last_seen"] = ""
	m["display_name"] = ""
	return m
}()

type authdata struct {
	client *Client

	// auth data from provider
	authobj   objx.Map
	hasCookie bool
}
