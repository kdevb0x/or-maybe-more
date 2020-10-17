// +import "github.com/kdevb0x/or-maybe-more/src/page"
package page

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

var (
	// root dir to serve static assets from (including html)
	StaticAssetDir string

	// paths needed for TLS.
	// for dev only
	// *NOT FOR PRODUCTION*
	devcert string
	devkey  string
)

var DefaultServer = new(Server)

type Server struct {
	http.Server
	Active bool
}

// updateHTML updates <htmlFile>.html by parsing the template name <htmlFile>.gohtml.
func (s *Server) UpdateHTML(htmlfile string) error {
	info, err := os.Stat(htmlfile)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("cannot parse %s, file doesn't seem to exist!", htmlfile))
		}
	}
	t, err := template.ParseFiles(htmlfile)
	if err != nil {
		throw(err, ERROR)
	}
	if wwwdir := os.Getenv("OMM_WWW_DIR"); wwwdir != "" && StaticAssetDir == "" {
		StaticAssetDir = wwwdir
	}

	fname := filepath.Join(StaticAssetDir, info.Name()+".html")
	f, err := os.Create(fname)
	if err != nil {
		// throw(err, ERROR)
		return err
	}
	if err := t.Execute(f, nil); err != nil {
		// throw(err, ERROR)
		return err
	}
	return nil
}

func ContactFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		throw(err, ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) Serve(addr string) {
	r := mux.NewRouter()
	sub := r.Host("www.ormaybemore.com").Subrouter()
	sub.Handle("/", http.FileServer(http.Dir(filepath.Join(StaticAssetDir, "www"))))
	s.Handler = r
	s.Active = true
	log.Fatal(s.ListenAndServeTLS(devcert, devkey))
}

// ContactInfo holds contact info submitted by a user to receive project
// announcments.
type ContactInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Tel       string `json:"telephone,omitempty"`
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

func addEmail(addr string, info ContactInfo) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	var info = ContactInfo
}
