package core

import (
	"crypto"
	"database/sql"

	_ "github.com/mattn/go-splite3"
)

type DBConn sql.DB

type DBAuthToken struct {
	PublicKey crypto.PublicKey
}

func OpenDB(dbpath string, auth DBAuthToken) (*DBConn, error) {
	h, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
}
