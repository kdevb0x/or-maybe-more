package page

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx"
	uuid "github.com/satori/go.uuid"
)

type DBConn struct {
	C *pgx.ConnPool
}

// Returns whether there is a database connection available from the pool.
func (db *DBConn) Connected() bool {
	//

	//
	return db.C.Stat().AvailableConnections > 0
}

// pgxDo is a dev helper to exexute pgx (postgres) operations without having to
// expose each method individually.
// It executes op using args, and returns the number of rows affected, and any
// errors.
func (db *DBConn) pgxDo(op string, args ...string) (int64, error) {
	var a = make([]interface{}, len(args))
	for i := range args {
		a[i] = args[i]
	}
	tag, err := db.C.Exec(op, a)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil

}

// QfSession is the QuikFace Aplication Programming Interface
type QfSession struct {
	URL      *url.URL
	Created  time.Time
	Duration time.Duration
	Key      crypto.PublicKey

	LastErr error

	// the number of active participants in the current sesssion.
	UserCount int

	// upper bound for participants; -1 == unlimited.
	UserLimit int

	// unexported local field context for the session
	ctx context.Context
}

func (s *Server) NewQfSession(ctx context.Context, userlimit int) *QfSession {
	ses := new(QfSession)
	ses.Created = time.Now()
	ses.UserLimit = userlimit
	if ctx != nil {
		ses.ctx = context.WithValue(ctx, nil, nil)
	}
	id := uuid.NewV4()
	// TODO: build url and make handler

	s.qfsessions = append(s.qfsessions, ses)
	s.QfSessions[id] = uint64(len(s.qfsessions) - 1)
	return ses
}

type Server struct {
	// embeds net/http.Server
	http.Server
	// contains page templates keyed by URL
	SiteDirectory map[string]uuid.UUID
	// cache of parsed templates keyed by their uuid
	cache map[uuid.UUID][]byte
	DB    DBConn

	// active quikface sessions
	qfsessions []*QfAPI

	// QFSessions contains the index into qfsessions, keyed by it's session
	// id (from which the url is generated).
	QfSessions map[uuid.UUID]uint64
}

// HtmlDocument fetches the rendered html of the document with uri if it exists,
// an error otherwise.
func (s *Server) HtmlDocument(uri string) ([]byte, error) {
	if id, ok := s.SiteDirectory[uri]; ok {
		return s.cache[id], nil
	}
	return nil, errors.New(fmt.Sprintf("document uri %s not found\n", uri))
}

// Handles www.ormaybemore.com/index/
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {

}

func CategoryHandler(w http.ResponseWriter, r *http.Request) {

}
