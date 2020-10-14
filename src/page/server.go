package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
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

var Server = new(server)

type server struct {
	http.Server
	active bool
}

// updateHTML updates <htmlFile>.html by parsing the template name <htmlFile>.gohtml.
func (s *server) updateHTML(htmlfile string) error {
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
	if wwwdir := os.Getenv("OMM_WWW_DIR"); wwwdir != "" {
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

func (s *server) serve(addr string) {
	r := mux.NewRouter()
	sub := r.Host("www.ormaybemore.com").Subrouter()
	sub.Handle("/", http.FileServer(http.Dir(filepath.Join(StaticAssetDir, "www"))))
	s.Handler = r
	s.active = true
	log.Fatal(s.ListenAndServeTLS(devcert, devkey))
}

func main() {
	errlog = make(chan taggedErr, 3)
	errlogger := new(errorLogger)

	// start processing loop
	go errlogger.run(os.Stderr)

	Server.updateHTML("index.gohtml")
	Server.serve(":8080")
}
