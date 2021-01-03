package page

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx"
	uuid "github.com/satori/go.uuid"
)

type Server struct {
	// embeds net/http.Server
	http.Server
	// contains page templates keyed by URL
	SiteDirectory map[string]uuid.UUID
	// cache of parsed templates keyed by their uuid
	cache  map[uuid.UUID][]byte
	DBConn *pgx.Conn
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
