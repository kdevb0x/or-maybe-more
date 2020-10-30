package net

import (
	"github.com/stretchr/objx"
)

type client struct {
}

type authdata struct {
	client *Client

	// auth data from provider
	authobj   objx.Map
	hasCookie bool
}
