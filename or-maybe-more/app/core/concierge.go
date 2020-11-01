package core

import (
	"crypto/tls"
	"time"
)

// AccessToken is used to manage room access.
type AccessToken tls.ClientSessionState

type room struct {
	displayName string
	uri         string
	created     time.Time
	available   bool
	manifest    *VacancyTable
}

// +gometatype singlton
//
// VacancyTable, (often called VTable for convenience) contains up to date information about the
// current room availability, overall capacity, and booking related logistics.
type VacancyTable struct {
	// IP address of the node responsible for issuing commands. Only one
	// node may have this responsibility at any given time.
	CommandNodeIP string
	rooms         []room
}

func (vt *VacancyTable) Claim(roomIdx int, claim ClaimToken) (AccessToken, error) {

}

func (vt *VacancyTable) Reserve(datetime time.Time) (*AccessToken, error) {

}
