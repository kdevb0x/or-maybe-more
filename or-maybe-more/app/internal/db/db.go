package core

import (
	"crypto"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	_ "modernc.org/sqlite"
)

type DBConn sql.DB

type DBAuthToken struct {
	PublicKey crypto.PublicKey
}

func OpenDB(dbpath string, auth DBAuthToken) (*DBConn, error) {
	h, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	if err != nil {
		return nil, err
	}
}

// ContactInfo holds contact info submitted by a user.
type ContactInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	TelNum    string `json:"telephone,omitempty"`
	Methods   ContactMethod
	Emails    []string `json:"email"`
}

func (ci *ContactInfo) Value() (driver.Value, error) {
	return json.Marshal(ci)
}

func (ci *ContactInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion of value to []byte failed")
	}
	return json.Unmarshal(b, &ci)
}

type ContactMethod uint

const (
	Text ContactMethod = 1 << iota
	Call
	Email
	// aka snailmail
	Letter
)
