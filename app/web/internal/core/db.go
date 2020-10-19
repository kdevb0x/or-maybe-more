package core

import (
	"crypto"
	"database/sql"

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
